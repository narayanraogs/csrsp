package db

import (
	"context"
	"csrsp/server/db/sqlc"
)

// Gets the max number of bands in the given frame type
func GetNoOfBands(frameType string, frameIdentifier string) (int64, error) {
	id, err := GetFrameIDForFrameIdentifier(frameType, frameIdentifier)
	if err != nil {
		return -1, err
	}
	var arg sqlc.CountOpticalSubTypesParams
	arg.Frameid = id
	arg.Subframetype = "Band"
	return global.CountOpticalSubTypes(context.Background(), arg)
}

// Gets the max number of ports in the given frame type
func GetNoOfPorts(frameType string, frameIdentifier string) (int64, error) {
	id, err := GetFrameIDForFrameIdentifier(frameType, frameIdentifier)
	if err != nil {
		return -1, err
	}
	var arg sqlc.CountOpticalSubTypesParams
	arg.Frameid = id
	arg.Subframetype = "Ports"
	return global.CountOpticalSubTypes(context.Background(), arg)
}

// Get the details of Optical Payload provided frame type and frame identifer
func GetOpticalFrameDetails(frameType string, frameIdentifier string) (sqlc.Videodataprocessing, error) {
	id, err := GetFrameIDForFrameIdentifier(frameType, frameIdentifier)
	if err != nil {
		return sqlc.Videodataprocessing{}, err
	}
	return global.GetVideoDataDetails(context.Background(), id)
}

// Get the Optical Details for all the Optical Payloads in the database
func GetAllOpticalDetails() (map[string]map[string]sqlc.Videodataprocessing, error) {
	var frameIdentifierMap = make(map[string]map[string]sqlc.Videodataprocessing)
	fts, err := GetAllOpticalFrameTypes()
	if err != nil {
		return nil, err
	}
	for _, frameType := range fts {
		frameIdentifierMap[frameType] = make(map[string]sqlc.Videodataprocessing)
		frames, err := GetFrameTypeFrameIdentifiers(frameType)
		if err != nil {
			return nil, err
		}
		for _, f := range frames {
			opt, err := global.GetVideoDataDetails(context.Background(), f.Frameid)
			if err != nil {
				return nil, err
			}
			frameIdentifierMap[frameType][f.Frameidentifier] = opt
		}
	}
	return frameIdentifierMap, nil
}

// Get the details of Optical Payload provided frame id
func GetOpticalFrameDetailsFrameIdentifer(frameID int) (sqlc.Videodataprocessing, error) {
	return global.GetVideoDataDetails(context.Background(), int32(frameID))
}

// Gets the Video limits for the provided FrameID
func GetVideoLimitsForFrameID(frameID int) (sqlc.Videolimit, error) {
	return global.GetOpticalLimitDetails(context.Background(), int32(frameID))
}

// Gets the Video limits for the provided Frame type and Frame Identifier
func GetVideoLimitsForFrameIdentifer(frameType string, frameID string) (sqlc.Videolimit, error) {
	fid, err := GetFrameIDForFrameIdentifier(frameType, frameID)
	if err != nil {
		return sqlc.Videolimit{}, err
	}
	return global.GetOpticalLimitDetails(context.Background(), fid)
}

// Gets the faulty pixels for the provided FrameID from the database
func GetFaultyPixels(frameID int) (string, error) {
	fp, err := global.GetOpticalFaultyPixels(context.Background(), int32(frameID))
	if err != nil {
		return "", err
	}
	return fp.Pixelnumbers, nil
}

// Get Sub frame details from Video SubFrames table
func GetSubFrames(subFrameType string, frameID int) ([]sqlc.Videosubframe, error) {
	var arg sqlc.GetOpticalSubFramesParams
	arg.Subframetype = sqlc.VideosubframesSubframetype(subFrameType)
	arg.Frameid = int32(frameID)
	return global.GetOpticalSubFrames(context.Background(), arg)
}

// Gets the subframe number for the provided sub frame name
func GetSubFrameNo(subFrameType string, frameID int, subFrameName string) (int, error) {
	var arg sqlc.GetOpticalSubFrameNumberParams
	arg.Subframetype = sqlc.VideosubframesSubframetype(subFrameType)
	arg.Frameid = int32(frameID)
	arg.Subframename = subFrameName
	sub, err := global.GetOpticalSubFrameNumber(context.Background(), arg)
	if err != nil {
		return -1, err
	}
	return int(sub.Subframenumber), nil
}

// GetStartEndPixelType ... return start, end and pixel type for the frame ID
func GetStartEndPixelType(subFrameType string, subFrameNo int, frameID int) (sqlc.Videosubframe, error) {
	var arg sqlc.GetOpticalSubFrameDetailsParams
	arg.Subframetype = sqlc.VideosubframesSubframetype(subFrameType)
	arg.Frameid = int32(frameID)
	arg.Subframenumber = int32(subFrameNo)
	return global.GetOpticalSubFrameDetails(context.Background(), arg)
}
