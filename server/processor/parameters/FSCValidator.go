// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrsp/server/db"
	"csrsp/server/session"
	"csrsp/server/utils/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strconv"
	"strings"
)

// FSCValidator checks a Frame Sequence Count (FSC) for bit errors.
type FSCValidator struct {
	paramName               string
	fsc                     []byte
	numberOfBitErrorAllowed int

	provider func(frame []byte) (string, bool)
}

// NewFSCValidator creates and-configures a new FSCValidator.
func NewFSCValidator(paramID string, paramName string) (*FSCValidator, error) {
	provider, err := getFSCValueProvider(paramID)
	if err != nil {
		return nil, err
	}

	fscDetails, err := db.GetFSCDetails(paramID)
	if err != nil {
		return nil, fmt.Errorf("no FSC details found for paramID %s", paramID)
	}

	tempArray := strings.Split(fscDetails.Fsc, ",")
	fsc, err := hex.DecodeString(strings.Join(tempArray, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to decode FSC hex string: %w", err)
	}

	return &FSCValidator{
		paramName:               paramName,
		fsc:                     fsc,
		numberOfBitErrorAllowed: int(fscDetails.Noofbiterrorallowed),
		provider:                provider,
	}, nil
}

// Process consumes an entire channel of frames, validating the FSC for each one.
func (v *FSCValidator) Process(input <-chan []byte, store *session.Store, key string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in FSCValidator.Process", "key", key, "panic", r, "stack", string(debug.Stack()))
		}
	}()

	var lineNo int64
	var errorCount int64
	var warningCount int64

	status, ok := store.Aux.GetStatus(key)
	if !ok {
		slog.Error("FSCValidator could not find aux status in store", "key", key)
		return
	}

	fscString := hex.EncodeToString(v.fsc)

	for array := range input {
		lineNo++
		fscGotString, ok := v.provider(array)
		if !ok {
			store.Aux.AddError(key, false, lineNo, fscString, fscGotString, "Corrupted Value in Data")
			status.SetError()
			errorCount++
			continue
		}

		fscGot, err := hex.DecodeString(fscGotString)
		if err != nil {
			store.Aux.AddError(key, false, lineNo, fscString, fscGotString, "Corrupted Value in Data")
			status.SetError()
			errorCount++
			continue
		}

		cardinality, err := binary.XORCardinality(v.fsc, fscGot)
		if err != nil {
			store.Aux.AddError(key, false, lineNo, fscString, fscGotString, "Unable to compute no of Bit Errors")
			status.SetError()
			errorCount++
		}
		if cardinality != 0 {
			if cardinality < v.numberOfBitErrorAllowed {
				store.Aux.AddError(key, true, lineNo, fscString, fscGotString, "No Of Bit Errors present, but below allowed limit")
				status.SetWarning()
				warningCount++
			} else {
				store.Aux.AddError(key, false, lineNo, fscString, fscGotString, "Bit Errors more than allowed limit")
				status.SetError()
				errorCount++
			}
		}
	}

	// Calculate and set final status
	var finalStatus string
	if errorCount == 0 && warningCount == 0 {
		finalStatus = "OK"
	} else if errorCount == 0 {
		finalStatus = "Warning [" + strconv.FormatInt(warningCount, 10) + " Times]"
	} else {
		finalStatus = "NOT OK [" + strconv.FormatInt(errorCount, 10) + " Errors]"
	}
	store.Aux.UpdateStatus(key, finalStatus)
}

func (v *FSCValidator) ParamName() string {
	return v.paramName
}
