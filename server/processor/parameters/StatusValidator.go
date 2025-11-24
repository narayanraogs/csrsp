// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrsp/server/db"
	"csrsp/server/session"
	"csrsp/server/utils/slice"
	"log/slog"
	"runtime/debug"
	"strconv"
	"strings"
)

// StatusValidator checks if a value corresponds to a set of expected status strings.
type StatusValidator struct {
	paramName             string
	multipleValuesAllowed bool
	provider              func([]byte) (uint64, string, bool)
}

func NewStatusValidator(paramID string, paramName string) (*StatusValidator, error) {
	provider, err := getStatusValueProvider(paramID)
	if err != nil {
		return nil, err
	}
	s, err := db.GetStatusMap(paramID)
	if err != nil {
		return nil, err
	}
	return &StatusValidator{
		paramName:             paramName,
		provider:              provider,
		multipleValuesAllowed: s.MultipleValueAllowed,
	}, nil
}

func (v *StatusValidator) Process(input <-chan []byte, store *session.Store, key string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in StatusValidator.Process", "key", key, "panic", r, "stack", string(debug.Stack()))
		}
	}()
	var prevValue string
	var first = true
	var lineNo int64
	var differentValues = make([]string, 0)
	var warningCount int64
	var errorCount int64

	stat, ok := store.Aux.GetStatus(key)
	if !ok {
		slog.Error("StatusValidator could not find aux status in store", "key", key)
		return
	}

	for array := range input {
		lineNo++
		value, statusGot, ok := v.provider(array)
		if !ok {
			store.Aux.AddError(v.paramName, false, lineNo, "",
				statusGot, "Not in Expected Values")
			stat.SetError()
			errorCount++
			continue
		}
		store.Aux.UpdateValue(key, strconv.FormatUint(value, 10))
		if first {
			prevValue = statusGot
			differentValues = append(differentValues, statusGot)
			first = false
			continue
		}
		if strings.Compare(strings.ToLower(prevValue), strings.ToLower(statusGot)) != 0 {
			if !v.multipleValuesAllowed {
				stat.SetWarning()
			}
			prevValue = statusGot
			if slice.IndexOfStringFold(differentValues, statusGot) == -1 {
				differentValues = append(differentValues, statusGot)
			}
			warningCount++
			continue
		}
	}
	var finalStatus string
	for _, value := range differentValues {
		finalStatus = finalStatus + value + ", "
	}
	finalStatus = strings.TrimSuffix(finalStatus, ", ")

	if errorCount > 0 {
		finalStatus = "NOT OK [" + strconv.FormatInt(errorCount, 10) + " Errors]"
	} else if len(differentValues) > 1 {
		finalStatus = "Multiple Values [" + finalStatus + "]"
	}

	store.Aux.UpdateStatus(key, finalStatus)
}

func (v *StatusValidator) ParamName() string {
	return v.paramName
}
