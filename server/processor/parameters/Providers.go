// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrsp/server/db"
	"csrsp/server/utils/binary"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/knetic/govaluate"
)

// --- Raw Value Providers ---

func getRawValueProvider(paramID string) (func(frame []byte) (uint64, bool), error) {
	bitMask, err := getBitMask(paramID)
	if err != nil {
		return nil, err
	}
	provider := func(frame []byte) (uint64, bool) {
		value, err := bitMask.ExtractUint64(frame)
		if err != nil {
			slog.Error("Cannot get Raw value for the Parameter", "paramID", paramID, "error", err)
			return 0, false
		}
		return value, true
	}
	return provider, nil
}

func getRawBitStringProvider(paramID string) (func(frame []byte) (string, bool), error) {
	bitMask, err := getBitMask(paramID)
	if err != nil {
		return nil, err
	}
	provider := func(frame []byte) (string, bool) {
		value, err := bitMask.ExtractString(frame)
		if err != nil {
			slog.Error("Cannot get Raw bit string for the Parameter", "paramID", paramID, "error", err)
			return "", false
		}
		return value, true
	}
	return provider, nil
}

func getBitMask(paramID string) (*binary.Mask, error) {
	param, err := db.GetParameterDetails(paramID)
	if err != nil {
		return nil, fmt.Errorf("parameter ID %s not found", paramID)
	}

	if param.Bitwise {
		bitWise, err := db.GetBitWiseDetailsForParamID(paramID)
		if err != nil {
			return nil, fmt.Errorf("ParamID %s not present in Bit Details Table", paramID)
		}
		var wordNo []int
		var validBits []string
		for _, bit := range bitWise {
			wordNo = append(wordNo, int(bit.Wordno))
			validBits = append(validBits, bit.Validbits)
		}
		return binary.NewBitwiseMask(wordNo, validBits)
	} else {
		noOfBits, err := strconv.Atoi(param.Noofbits)
		if err != nil {
			return nil, fmt.Errorf("ParamID %s No of Bits is not correct %d", paramID, noOfBits)
		}
		return binary.NewContinuousMask(int(param.Startword), int(param.Startbit), noOfBits)
	}
}

// --- Processed Value Providers ---

func getStatusValueProvider(paramID string) (func([]byte) (uint64, string, bool), error) {
	rawProvider, err := getRawValueProvider(paramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw value provider for %s: %w", paramID, err)
	}
	statusMap, err := db.GetStatusMap(paramID)
	if err != nil {
		return nil, fmt.Errorf("cannot get status map for paramID %s", paramID)
	}

	provider := func(frame []byte) (uint64, string, bool) {
		value, ok := rawProvider(frame)
		if !ok {
			return 0, "", false
		}
		interpretation, ok := statusMap.Status[value]
		if !ok {
			return value, strconv.FormatUint(value, 10), false
		}
		return value, interpretation, true
	}
	return provider, nil
}

func getFSCValueProvider(paramID string) (func(frame []byte) (string, bool), error) {
	raw, err := getRawBitStringProvider(paramID)
	if err != nil {
		return nil, fmt.Errorf("unable to get RawValueprovider for ParamID %s: %w", paramID, err)
	}

	provider := func(frame []byte) (string, bool) {
		bitString, ok := raw(frame)
		if !ok {
			return "", false
		}
		value, err := binary.BitStringToHexString(bitString)
		if err != nil {
			return "", false
		}
		return value, true
	}
	return provider, nil
}

func getIncrementValueProvider(paramID string) (func(frame []byte) (uint64, bool), error) {
	return getRawValueProvider(paramID)
}

func getRadixValueProvider(paramID string) (func(frame []byte) (uint64, string, bool), error) {
	rawProvider, err := getRawValueProvider(paramID)
	if err != nil {
		return nil, err
	}
	var radix int
	r, err := db.GetRadixDetails(paramID)
	if err != nil {
		return nil, fmt.Errorf("radixTable entry is not proper for ParamID %s", paramID)
	}
	if r.Radix <= 0 {
		radix = 10
	}

	provider := func(frame []byte) (uint64, string, bool) {
		value, ok := rawProvider(frame)
		if !ok {
			return 0, "", false
		}
		return value, strconv.FormatUint(value, radix), true
	}
	return provider, nil
}

func getCRCValueProvider(paramID string) (func(frame []byte) (uint16, bool), error) {
	rawProvider, err := getRawValueProvider(paramID)
	if err != nil {
		return nil, err
	}

	provider := func(frame []byte) (uint16, bool) {
		value, ok := rawProvider(frame)
		if !ok {
			return 0, false
		}
		return uint16(value), true
	}
	return provider, nil
}

func getAnalogValueProvider(paramID string) (func([]byte) (float64, bool), error) {
	analog, err := db.GetAnalogDetails(paramID)
	if err != nil {
		return nil, err
	}
	variable := strings.TrimSpace(analog.Variablename)

	calculatedValues := make(map[string]float64)
	parameters := make(map[string]interface{})

	paramIDs, variables := getParameterIDs(analog.Parameterids.String)
	providers := make([]func([]byte) interface{}, len(paramIDs))
	for i := 0; i < len(paramIDs); i++ {
		providers[i], err = GetProcessedParameterValueProvider(paramIDs[i])
		if err != nil {
			providers[i] = nil
		}
	}
	paramIDs = append(paramIDs, paramID)
	variables = append(variables, variable)
	providers = append(providers, getParameterValueProvider(paramID, string(analog.Datatype)))

	exp, err := govaluate.NewEvaluableExpression(analog.Equation)
	if err != nil {
		return nil, err
	}
	var provider = func(frame []byte) (float64, bool) {
		var calcValue string
		for i := 0; i < len(paramIDs); i++ {
			value := providers[i](frame)
			parameters[variables[i]] = value
			calcValue = calcValue + fmt.Sprintf("%v", value) + ";"
		}
		analogValue, ok := calculatedValues[calcValue]
		if ok {
			return analogValue, true
		}
		result, err := exp.Evaluate(parameters)
		if err != nil {
			return 0.0, false
		}
		return result.(float64), true
	}
	return provider, nil
}

func getParameterIDs(parameterIDs string) ([]string, []string) {
	param := make([]string, 0)
	variables := make([]string, 0)
	parameterIDs = strings.TrimSpace(parameterIDs)
	if strings.EqualFold(parameterIDs, "") {
		return param, variables
	}
	temp := strings.Split(parameterIDs, ",")
	for _, parameter := range temp {
		variables = append(variables, strings.TrimSpace(parameter))
		parameter = strings.ToUpper(parameter)
		parameter = strings.ReplaceAll(parameter, "P", "")
		id := strings.TrimSpace(parameter)
		param = append(param, id)

	}
	return param, variables
}

func GetProcessedParameterValueProvider(paramID string) (func([]byte) interface{}, error) {
	param, err := db.GetParameterDetails(paramID)
	if err != nil {
		return nil, err
	}
	raw, err := getRawValueProvider(paramID)
	if err != nil {
		return nil, err
	}
	rawProvider := func(frame []byte) interface{} {
		value, valueOK := raw(frame)
		if !valueOK {
			slog.Error("Cannot process parameter in Analog Value Provider", "ParamID", paramID)
			return 0
		}
		return value
	}

	switch strings.ToLower(string(param.Parametertype)) {
	case "status":
		statusProvider, err := getStatusValueProvider(paramID)
		if err != nil {
			return nil, err
		}
		provider := func(frame []byte) interface{} {
			_, value, statusOK := statusProvider(frame)
			if !statusOK {
				slog.Error("Cannot process parameter in status Value Provider", "ParamID", paramID)
				return ""
			}
			return value
		}
		return provider, nil
	case "radix":
		return rawProvider, nil
	case "increment":
		return rawProvider, nil
	case "fsc":
		return rawProvider, nil
	case "analog":
		analog, err := getAnalogValueProvider(paramID)
		if err != nil {
			return nil, err
		}
		analogProvider := func(frame []byte) interface{} {
			analogValue, analogOK := analog(frame)
			if !analogOK {
				slog.Error("Cannot process parameter in Analog Value Provider", "ParamID", paramID)
				return 0
			}
			return analogValue
		}
		return analogProvider, nil
	default:
		return rawProvider, nil
	}
}

func getParameterValueProvider(paramID string, dataType string) func([]byte) interface{} {
	raw, err := getRawValueProvider(paramID)
	if err != nil {
		slog.Error("Cannot get Raw value provider ", "ParamID", paramID)
		return nil
	}
	provider := func(frame []byte) interface{} {
		value, ok := raw(frame)
		if !ok {
			slog.Error("Cannot read value for paramter", "ParamID", paramID)
			return 0
		}
		switch strings.ToLower(dataType) {
		case "int_16":
			return int16(value)
		case "int_32":
			return int32(value)
		case "int_64":
			return int64(value)
		case "uint_16":
			return value
		case "uint_32":
			return value
		case "uint_64":
			return value
		case "float_ieee_32":
			return math.Float32frombits(uint32(value & 0xFFFFFFFF))
		case "float_ieee_64":
			return math.Float64frombits(value)
		case "float_1750a_32":
			return binary.DecodeF1750A32(value)
		case "float_1750a_48":
			return binary.DecodeF1750A48(value)
		default:
			return value
		}
	}
	return provider
}
