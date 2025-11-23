// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrspServer/db"
	"csrspServer/session"
	"csrspServer/utils/slice"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strconv"
	"strings"
)

// AnalogValidator checks if a numeric value and its difference from the previous
// value are within specified limits.
type AnalogValidator struct {
	paramName string

	// Config
	valueValidationRequired bool
	diffValidationRequired  bool
	upperLimit              float64
	lowerLimit              float64
	diffUpperLimit          float64
	diffLowerLimit          float64

	provider func([]byte) (float64, bool)
}

// NewAnalogValidator creates and configures a new AnalogValidator.
func NewAnalogValidator(paramID, paramName string) (*AnalogValidator, error) {
	provider, err := getAnalogValueProvider(paramID) // Assumes getAnalogValueProvider is ported
	if err != nil {
		return nil, err
	}

	analog, err := db.GetAnalogDetails(paramID)
	if err != nil {
		return nil, fmt.Errorf("no Analog details found for paramID %s", paramID)
	}

	return &AnalogValidator{
		paramName:               paramName,
		provider:                provider,
		valueValidationRequired: analog.Valuevalidationrequired,
		diffValidationRequired:  analog.Differencevalidationrequired,
		upperLimit:              analog.Upperlimitvalue.Float64 + analog.Tolerancevalue.Float64,
		lowerLimit:              analog.Lowerlimitvalue.Float64 - analog.Tolerancevalue.Float64,
		diffUpperLimit:          analog.Differencevalue.Float64 + analog.Differencetolerance.Float64,
		diffLowerLimit:          analog.Differencevalue.Float64 - analog.Differencetolerance.Float64,
	}, nil
}

// Process consumes an entire channel of frames, validating the analog value for each one.
func (v *AnalogValidator) Process(input <-chan []byte, store *session.Store, key string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in AnalogValidator.Process", "key", key, "panic", r, "stack", string(debug.Stack()))
		}
	}()

	var lineNo int64
	var first = true
	var prevValue, minValue, maxValue float64
	var differentValues []string

	stat, ok := store.Aux.GetStatus(key)
	if !ok {
		slog.Error("AnalogValidator could not find aux status in store", "key", key)
		return
	}

	for array := range input {
		lineNo++
		value, ok := v.provider(array)
		strValue := strconv.FormatFloat(value, 'G', -1, 64)
		if !ok {
			store.Aux.AddError(v.paramName, false, lineNo, "",
				strValue, "Corrupted value in data")
			stat.SetError()
			continue
		}
		store.Aux.UpdateValue(key, strValue)

		if v.valueValidationRequired {
			if value < v.lowerLimit || value > v.upperLimit {
				store.Aux.AddError(v.paramName, false, lineNo, "",
					strValue, "Value not within provided limits")
				stat.SetError()
				continue
			}
		}

		if len(differentValues) <= 10 && slice.IndexOfStringFold(differentValues, strValue) == -1 {
			differentValues = append(differentValues, strValue)
		}

		if first {
			first = false
			prevValue = value
			minValue = value
			maxValue = value
			continue
		}

		if value < minValue {
			minValue = value
		}

		if value > maxValue {
			maxValue = value
		}

		if v.diffValidationRequired {
			var diff float64
			if prevValue > value {
				diff = prevValue - value
			} else {
				diff = value - prevValue
			}
			if diff < v.diffLowerLimit || diff > v.diffUpperLimit {
				store.Aux.AddError(v.paramName, false, lineNo, "",
					strValue, "Difference not within provided limits")
				stat.SetError()
			}
		}
		prevValue = value
	}

	// Calculate and set final status
	var finalStatus string
	if len(differentValues) <= 10 {
		finalStatus = strings.Join(differentValues, ",")
	} else {
		finalStatus = fmt.Sprintf("%f to %f", minValue, maxValue)
	}
	store.Aux.UpdateStatus(key, finalStatus)
}

func (v *AnalogValidator) ParamName() string {
	return v.paramName
}
