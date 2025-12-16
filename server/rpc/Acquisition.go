package rpc

import (
	"context"
	pb "csrsp/server/communication"
	"csrsp/server/das"
	"csrsp/server/db"
	"fmt"
	"log"
	"log/slog"
	"time"
)

func (s *CommunicationServer) GetAcquisitionParameters(ctx context.Context, req *pb.ClientID) (*pb.AcquisitionParameters, error) {
	var acq pb.AcquisitionParameters
	acq.AcqTypes = []string{"Frame", "Time", "User Controlled"}
	acqModes, err := db.GetAcquisitionModes("Acquisition")
	if err != nil {
		return nil, err
	}
	acq.AcqModes = acqModes
	configNames, err := db.GetAllConfigNames()
	if err != nil {
		return nil, err
	}
	acq.ConfigNames = configNames
	payloads, err := db.GetAllPayloadNames()
	if err != nil {
		return nil, err
	}
	acq.Payloads = payloads
	resultProfiles, err := db.GetAllResultProfiles()
	if err != nil {
		return nil, err
	}
	acq.ResultProfiles = resultProfiles
	acq.DasMap = make([]*pb.AcqDASMap, 0)
	for _, acqMode := range acqModes {
		var d pb.AcqDASMap
		d.AcqMode = acqMode
		d.DasDetails = make([]*pb.AcqDasDetails, 0)
		dasPaths, err := db.GetAcquisitionChainDetails(acqMode)
		if err != nil {
			return nil, err
		}
		for _, path := range dasPaths {
			var dd pb.AcqDasDetails
			dd.DpuNumber = path.Dpunumber

			// Get DAS Name
			dasInfo, err := db.GetDASDetails(int(path.Dasid))
			if err != nil {
				return nil, err
			}
			dd.DasName = dasInfo.Dasname

			d.DasDetails = append(d.DasDetails, &dd)
		}
		acq.DasMap = append(acq.DasMap, &d)
	}

	return &acq, nil
}

func (s *CommunicationServer) GetDASStatus(req *pb.DASStatusRequest, stream pb.Communication_GetDASStatusServer) error {
	log.Printf("Starting das stream for client: %s", req.Id)

	// Stream data every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Helper function to fetch and send status
	sendUpdate := func() error {
		status, err := das.GetStatus(req.AcqMode)
		if err != nil {
			slog.Error("Unable to get DAS Status", "error", err.Error())
			return nil // Don't break stream on calculation error, just skip sending
		}

		var dasStatus []*pb.DASStatusResponse
		for _, s := range status {
			dasStatus = append(dasStatus, &pb.DASStatusResponse{
				DasName:   s.DasName,
				DpuNumber: int32(s.DpuNumber),
				Status:    s.Status,
				Alarm:     s.Alarm,
			})
		}

		return stream.Send(&pb.DASStatus{DasStatus: dasStatus})
	}

	// 1. Send immediate update
	if err := sendUpdate(); err != nil {
		return fmt.Errorf("failed to send initial status: %v", err)
	}

	// 2. Loop for periodic updates
	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Client %s disconnected", req.Id)
			return nil
		case <-ticker.C:
			if err := sendUpdate(); err != nil {
				return fmt.Errorf("failed to send status: %v", err)
			}
		}
	}
}

func (s *CommunicationServer) GetFileAcquisitionParameters(ctx context.Context, req *pb.ClientID) (*pb.FileAcquisitionParameters, error) {
	var params pb.FileAcquisitionParameters
	acqModes, err := db.GetAcquisitionModes("Acquisition")
	if err != nil {
		return nil, err
	}
	params.AcqModes = acqModes
	configNames, err := db.GetAllConfigNames()
	if err != nil {
		return nil, err
	}
	params.ConfigNames = configNames
	payloads, err := db.GetAllPayloadNames()
	if err != nil {
		return nil, err
	}
	params.Payloads = payloads
	resultProfiles, err := db.GetAllResultProfiles()
	if err != nil {
		return nil, err
	}
	params.ResultProfiles = resultProfiles
	frameTypes, err := db.GetAllFrameTypes()
	if err != nil {
		return nil, err
	}
	params.FrameTypes = frameTypes
	params.FrameTypeMap = make([]*pb.FrameTypeMap, 0)
	for _, frameType := range frameTypes {
		fids := make([]string, 0)
		frameIdentifiers, err := db.GetFrameTypeFrameIdentifiers(frameType)
		if err != nil {
			return nil, err
		}
		for _, fid := range frameIdentifiers {
			fids = append(fids, fid.Frameidentifier)
		}
		params.FrameTypeMap = append(params.FrameTypeMap, &pb.FrameTypeMap{
			FrameType:        frameType,
			FrameIdentifiers: fids,
		})
	}
	return &params, nil
}
