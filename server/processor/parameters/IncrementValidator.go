// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrspServer/db"
	"csrspServer/session"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strconv"
)

// IncrementValidator checks if a value increments correctly.
type IncrementValidator struct {
	paramName  string
	increment  int
	startValue int
	endValue   int
	tolerance  int
	maxValue   uint64
	provider   func(frame []byte) (uint64, bool)
}

func NewIncrementValidator(paramID string, paramName string) (*IncrementValidator, error) {
	provider, err := getIncrementValueProvider(paramID)
	if err != nil {
		return nil, err
	}
	incr, err := db.GetIncrementDetails(paramID)
	if err != nil {
		return nil, err
	}
	p, err := db.GetParameterDetails(paramID)
	if err != nil {
		return nil, err
	}
	noOfBits, err := strconv.Atoi(p.Noofbits)
	if err != nil {
		return nil, err
	}
	return &IncrementValidator{
		paramName:  paramName,
		startValue: int(incr.Startvalue),
		endValue:   int(incr.Endvalue),
		tolerance:  int(incr.Tolerance),
		increment:  int(incr.Increment),
		maxValue:   1 << noOfBits,
		provider:   provider,
	}, nil
}

func (v *IncrementValidator) Process(input <-chan []byte, store *session.Store, key string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in IncrementValidator.Process", "key", key, "panic", r, "stack", string(debug.Stack()))
		}
	}()
	var first = true
	var lineNo int64
	var valuesPresent bool
	var prevValue uint64
	stat, ok := store.Aux.GetStatus(key)
	if !ok {
		slog.Error("AnalogValidator could not find aux status in store", "key", key)
		return
	}
	var errorCount int64

	if v.startValue == -1 || v.endValue == -1 || v.startValue == v.endValue {
		valuesPresent = false
	} else {
		valuesPresent = true
	}

	for array := range input {
		value, ok := v.provider(array)
		lineNo++
		if !ok {
			store.Aux.AddError(v.paramName, false, lineNo, "",
				strconv.Itoa(int(value)), "Corrupted Value in Data")
			stat.SetError()
			errorCount++
			continue
		}
		store.Aux.UpdateValue(key, strconv.FormatUint(value, 10))
		if first {
			first = false
			prevValue = value
			continue
		}
		var diff uint64
		expectedValue := prevValue + uint64(v.increment)
		if valuesPresent {
			expectedValue = expectedValue % (uint64(v.endValue) + 1)
		}
		expectedValue = expectedValue % v.maxValue
		if value < expectedValue {
			diff = expectedValue - value
		} else {
			diff = value - expectedValue
		}
		if diff > uint64(v.tolerance) {
			store.Aux.AddError(v.paramName, false, lineNo, strconv.FormatUint(expectedValue, 10),
				strconv.FormatUint(value, 10), "Difference not within provided limits")
			stat.SetError()
			prevValue = value
			errorCount++
			continue
		}
		if valuesPresent && value < uint64(v.startValue) {
			store.Aux.AddError(v.paramName, false, lineNo, fmt.Sprintf("%d", v.startValue),
				strconv.FormatUint(value, 10), "Received value less than Start Value")
			stat.SetError()
			prevValue = value
			errorCount++
			continue
		}
		if valuesPresent && value > uint64(v.endValue) {
			store.Aux.AddError(v.paramName, false, lineNo, fmt.Sprintf("%d", v.endValue),
				strconv.FormatUint(value, 10), "Received value more than End Value")
			stat.SetError()
			prevValue = value
			errorCount++
			continue
		}
		prevValue = value
	}

	var finalStatus string
	if errorCount == 0 {
		finalStatus = "OK"
	} else {
		finalStatus = "NOT-OK [" + strconv.FormatInt(errorCount, 10) + " Errors]"
	}
	store.Aux.UpdateStatus(key, finalStatus)
}

func (v *IncrementValidator) ParamName() string {
	return v.paramName
}
