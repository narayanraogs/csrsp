package db

import (
	"context"
	"csrsp/server/db/sqlc"
	"strconv"
	"strings"
)

type StatusMap struct {
	ParamID              string
	Status               map[uint64]string
	MultipleValueAllowed bool
}

// Gets the parameter details provided the Parameter ID
func GetParameterDetails(paramID string) (sqlc.Parameterdetail, error) {
	return global.GetParameterDetails(context.Background(), paramID)
}

// Gets the Parameter ID provide the frame type and parameter name
func GetPrameterID(frameType string, parameterName string) (string, error) {
	var args sqlc.GetParamIDParams
	args.Frametype = frameType
	args.Parametername = parameterName
	return global.GetParamID(context.Background(), args)
}

// Gets the Bit wise details for the provided parameter ID
func GetBitWiseDetailsForParamID(paramID string) ([]sqlc.Bitdetail, error) {
	return global.GetBitDetails(context.Background(), paramID)
}

// Gets the FSC Details for the provided Parameter ID
func GetFSCDetails(paramID string) (sqlc.Fsctype, error) {
	return global.GetFSCDetails(context.Background(), paramID)
}

// Gets the CRC Details for the provided Parameter ID
func GetCRCDetails(paramID string) (sqlc.Crctype, error) {
	return global.GetCRCDetails(context.Background(), paramID)
}

// Gets the Increment Details for the provided Parameter ID
func GetIncrementDetails(paramID string) (sqlc.Incrementtype, error) {
	return global.GetIncrementDetails(context.Background(), paramID)
}

// Gets the Radix Details for the provided Parameter ID
func GetRadixDetails(paramID string) (sqlc.Radixtype, error) {
	return global.GetRadixDetails(context.Background(), paramID)
}

// Gets the Analog Details for the provided Parameter ID
func GetAnalogDetails(paramID string) (sqlc.Analogtype, error) {
	return global.GetAnalogDetails(context.Background(), paramID)
}

// Lists all the parameters for the provied frame type
func GetFrameTypeParamterDetails(frameType string) ([]sqlc.Parameterdetail, error) {
	return global.GetFrameTypeParamterDetails(context.Background(), frameType)
}

// Gets the status map for the specified parameter Id
func GetStatusMap(paramID string) (StatusMap, error) {
	values, err := global.GetStatusDetails(context.Background(), paramID)
	if err != nil {
		return StatusMap{}, err
	}
	var statusMap StatusMap
	statusMap.ParamID = paramID
	statusMap.MultipleValueAllowed = values[0].Multiplevaluesallowed
	statusMap.Status = make(map[uint64]string)
	for _, stat := range values {
		base := "0d"
		bitString := stat.Bitvalue
		if len(stat.Bitvalue) > 2 {
			base = strings.ToLower(strings.TrimSpace(stat.Bitvalue[0:2]))
			bitString = stat.Bitvalue[2:]
		}
		var value uint64
		switch base {
		case "0x":
			value, _ = strconv.ParseUint(bitString, 16, 64)
		case "0b":
			value, _ = strconv.ParseUint(bitString, 2, 64)
		case "0d":
			value, _ = strconv.ParseUint(bitString, 10, 64)
		default:
			value, _ = strconv.ParseUint(stat.Bitvalue, 10, 64)
		}
		statusMap.Status[value] = stat.Interpretation
	}
	return statusMap, nil
}

// Gets the FSC value as a string for the provided frameID
func GetFSCValueForFrameID(frameID int) (string, error) {
	frameType, err := GetFrameTypeDetailsForFrameID(frameID)
	if err != nil {
		return "", err
	}

	pid, err := global.GetFSCParamID(context.Background(), frameType.Frametype)
	if err != nil {
		return "", err
	}

	fsc, err := GetFSCDetails(pid)
	if err != nil {
		return "", err
	}
	return fsc.Fsc, nil
}
