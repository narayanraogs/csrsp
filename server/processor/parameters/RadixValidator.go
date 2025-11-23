// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrspServer/db"
	"csrspServer/session"
	"csrspServer/utils/slice"
	"log/slog"
	"runtime/debug"
	"strconv"
	"strings"
)

// RadixValidator checks a value against configured limits.
type RadixValidator struct {
	paramName                    string
	upperLimit                   uint64
	lowerLimit                   uint64
	diffUpperLimit               int64
	diffLowerLimit               int64
	valueValidationRequired      bool
	differenceValidationRequired bool
	provider                     func(frame []byte) (uint64, string, bool)
}

func NewRadixValidator(paramID string, paramName string) (*RadixValidator, error) {
	provider, err := getRadixValueProvider(paramID)
	if err != nil {
		return nil, err
	}
	r, err := db.GetRadixDetails(paramID)
	if err != nil {
		return nil, err
	}
	return &RadixValidator{
		paramName:                    paramName,
		provider:                     provider,
		valueValidationRequired:      r.Valuevalidationrequired,
		upperLimit:                   uint64(r.Upperlimitvalue.Int32 + r.Valuetolerance.Int32),
		lowerLimit:                   uint64(r.Upperlimitvalue.Int32 - r.Valuetolerance.Int32),
		differenceValidationRequired: r.Differencevalidationrequired,
		diffUpperLimit:               int64(r.Differencevalue.Int32 + r.Differencetolerance.Int32),
		diffLowerLimit:               int64(r.Differencevalue.Int32 - r.Differencetolerance.Int32),
	}, nil
}

func (v *RadixValidator) Process(input <-chan []byte, store *session.Store, key string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in RadixValidator.Process", "key", key, "panic", r, "stack", string(debug.Stack()))
		}
	}()

	var prevValue uint64
	var first = true
	var lineNo int64
	var minValue uint64
	var maxValue uint64
	var differentValues []string
	differentValues = make([]string, 0)

	stat, ok := store.Aux.GetStatus(key)
	if !ok {
		slog.Error("AnalogValidator could not find aux status in store", "key", key)
		return
	}

	for array := range input {
		lineNo++
		value, strValue, ok := v.provider(array)
		if !ok {
			store.Aux.AddError(v.paramName, false, lineNo, "",
				strconv.FormatUint(value, 10), "Corrupted value in data")
			stat.SetError()
			continue
		}
		store.Aux.UpdateValue(key, strconv.FormatUint(value, 10))
		if len(differentValues) <= 10 && slice.IndexOfStringFold(differentValues, strValue) == -1 {
			differentValues = append(differentValues, strValue)
		}

		if v.valueValidationRequired {
			if value < v.lowerLimit || value > v.upperLimit {
				store.Aux.AddError(v.paramName, false, lineNo, "",
					strValue, "Value not within provided limits")
				stat.SetError()
				continue
			}
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

		if v.differenceValidationRequired {
			var diff uint64
			if prevValue > value {
				diff = prevValue - value
			} else {
				diff = value - prevValue
			}
			if diff < uint64(v.diffLowerLimit) || diff > uint64(v.diffUpperLimit) {
				store.Aux.AddError(v.paramName, false, lineNo, "",
					strValue, "Difference not within provided limits")
				stat.SetError()
			}
		}
		prevValue = value
	}

	var finalStatus string
	if len(differentValues) <= 10 {
		for _, value := range differentValues {
			finalStatus = finalStatus + value + ", "
		}
		finalStatus = strings.TrimSuffix(finalStatus, ", ")
	} else {
		finalStatus = strconv.FormatUint(minValue, 10)
		finalStatus += " to "
		finalStatus += strconv.FormatUint(maxValue, 10)
	}
	store.Aux.UpdateStatus(key, finalStatus)
}

func (v *RadixValidator) ParamName() string {
	return v.paramName
}
