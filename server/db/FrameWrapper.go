package db

import (
	"context"
	"csrspServer/db/sqlc"
	"slices"
)

// Gets the frame details for the specified frame ID
func GetFrameDetails(frameID int) (sqlc.Frame, error) {
	return global.GetFrameDetails(context.Background(), int32(frameID))
}

// Gets all the Frames within the Specificed Frame type
func GetFrameTypeFrameIdentifiers(frameType string) ([]sqlc.Frame, error) {
	return global.GetFrameTypeFrameIdentifiers(context.Background(), frameType)
}

// Gets all the Frame types to be archived
func GetArchivableFrameTypes() ([]string, error) {
	return global.GetArchivableFrameTypes(context.Background())
}

// Checks if the Provided frame type is to be archived
func IsFrameTypeToBeArchived(frameType string) (bool, error) {
	values, err := global.GetArchivableFrameTypes(context.Background())
	if err != nil {
		return false, err
	}
	return slices.Contains(values, frameType), nil
}

// Checks if the Provided frame type is to be written
func IsFrameTypeToBeWritten(frameType string) (bool, error) {
	values, err := global.GetWritableFrameTypes(context.Background())
	if err != nil {
		return false, err
	}
	return slices.Contains(values, frameType), nil
}

// Checks if the Provided frame type is an acquired frame type
func IsFrameTypeAcquired(frameType string) (bool, error) {
	values, err := global.GetAcquiredFrameTypes(context.Background())
	if err != nil {
		return false, err
	}
	return slices.Contains(values, frameType), nil
}

// Gets all the Aux Frame Types
func GetAllAuxFrameTypes() ([]string, error) {
	return global.GetAuxFrameTypes(context.Background())
}

// Gets all the Opitcal Frame Types
func GetAllOpticalFrameTypes() ([]string, error) {
	return global.GetOpticalFrameTypes(context.Background())
}

// Gets all the Opitcal Frame Types
func GetAllFrameTypes() ([]string, error) {
	return global.GetFrameTypes(context.Background())
}

// Gets the frame id for the frame type and frame identifier
func GetFrameIDForFrameIdentifier(frameType string, frameIdentifier string) (int32, error) {
	var arg sqlc.GetFrameIDParams
	arg.Frameidentifier = frameIdentifier
	arg.Frametype = frameType
	return global.GetFrameID(context.Background(), arg)
}

// Gets all the Microwave Frame Types
func GetAllMicrowaveFrameTypes() ([]string, error) {
	return global.GetMicrowaveFrameTypes(context.Background())
}

// Gets the frame details for the provided frame type and frame identifer
func GetFrameDetailsForFrameIdentifer(frameType string, frameIdentifier string) (sqlc.Frame, error) {
	id, err := GetFrameIDForFrameIdentifier(frameType, frameIdentifier)
	if err != nil {
		return sqlc.Frame{}, err
	}
	return GetFrameDetails(int(id))
}

// Gets The type of the provided frame type
func GetTypeOfFrameType(frameType string) (string, error) {
	t, err := global.GetTypeOfFrame(context.Background(), frameType)
	if err != nil {
		return "", err
	}
	return string(t), nil
}

// Gets the Aux Intermediate Frame Types
func GetAuxIntermediateFrameTypes() ([]string, error) {
	return global.GetIntermediateFrames(context.Background(), "Aux")
}

// Gets the Optical Intermediate Frame Types
func GetOpticalIntermediateFrameTypes() ([]string, error) {
	return global.GetIntermediateFrames(context.Background(), "Optical")
}

// Gets all frame ids
func GetAllFrameIDs() ([]int32, error) {
	return global.GetAllFrameIDs(context.Background())
}

// Gets all the frames of the frame types present in the database
func GetAllFrames() ([]sqlc.Frame, error) {
	var tbr = make([]sqlc.Frame, 0)
	fts, err := GetAllFrameTypes()
	if err != nil {
		return nil, err
	}
	for _, ft := range fts {
		f, err := GetFrameTypeFrameIdentifiers(ft)
		if err != nil {
			return nil, err
		}
		tbr = append(tbr, f...)
	}
	return tbr, nil
}

// Gets the Frame type details for the specified frametype
func GetFrameTypeDetails(frameType string) (sqlc.Frametype, error) {
	return global.GetFrameTypeDetails(context.Background(), frameType)
}

// Gets the identifying param id for the provided frame id
func GetFrameTypeDetailsForFrameID(frameID int) (sqlc.Frametype, error) {
	frame, err := GetFrameDetails(frameID)
	if err != nil {
		return sqlc.Frametype{}, err
	}
	return GetFrameTypeDetails(frame.Frametype)
}
