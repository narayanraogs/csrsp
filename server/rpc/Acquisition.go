package rpc

import (
	"context"
	"csrsp/server/communication"
	"csrsp/server/db"
)

func (s *CommunicationServer) GetAcquisitionParameters(ctx context.Context, req *communication.ClientID) (*communication.AcquisitionParameters, error) {
	var acq communication.AcquisitionParameters
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
	acq.DasMap = make([]*communication.AcqDASMap, 0)
	for _, acqMode := range acqModes {
		var d communication.AcqDASMap
		d.AcqMode = acqMode
		d.DasDetails = make([]*communication.AcqDasDetails, 0)
		dasDetail, err := db.GetDASSystemsByAcquisitionMode(acqMode)
		if err != nil {
			return nil, err
		}
		for _, das := range dasDetail {
			var dd communication.AcqDasDetails
			dd.DasName = das.Dasname
			//dd.DpuNumber = das.Dpunumber
			d.DasDetails = append(d.DasDetails, &dd)
		}
		acq.DasMap = append(acq.DasMap, &d)
	}

	return &acq, nil
}
