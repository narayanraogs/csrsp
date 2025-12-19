package rpc

import (
	"context"
	pb "csrsp/server/communication"
	"csrsp/server/db"
)

func (s *CommunicationServer) GetAllTestPhases(ctx context.Context, req *pb.ClientID) (*pb.TestPhases, error) {
	var tps pb.TestPhases
	tps.TestPhases = make([]string, 0)
	var err error
	tps.TestPhases, err = db.GetTestPhases()
	return &tps, err
}

func (s *CommunicationServer) AddTestPhase(ctx context.Context, req *pb.TestPhaseRequest) (*pb.Ack, error) {
	var ack pb.Ack
	err := db.CreateTestPhase(req.TestPhase)
	if err != nil {
		ack.Ok = false
		ack.Message = "Failed to add test phase"
		return &ack, err
	}
	ack.Ok = true
	ack.Message = "Test phase added successfully"
	return &ack, err
}

func (s *CommunicationServer) SelectTestPhase(ctx context.Context, req *pb.TestPhaseRequest) (*pb.Ack, error) {
	var ack pb.Ack
	err := db.ChangeTestPhase(req.TestPhase)
	if err != nil {
		ack.Ok = false
		ack.Message = "Failed to change test phase"
		return &ack, err
	}
	ack.Ok = true
	ack.Message = "Test phase changed successfully"
	return &ack, err
}

func (s *CommunicationServer) GetDASIPAddresses(ctx context.Context, req *pb.ClientID) (*pb.DASIPAddresses, error) {
	var dasIPAddresses pb.DASIPAddresses
	dasIPAddresses.DasIPAddresses = make([]*pb.DASIPAddress, 0)
	dasDetail, err := db.GetAllDASDetails()
	if err != nil {
		return nil, err
	}
	for _, das := range dasDetail {
		var dasIPAddress pb.DASIPAddress
		dasIPAddress.Name = das.Dasname
		dasIPAddress.IpAddress = das.Ipaddress
		dasIPAddresses.DasIPAddresses = append(dasIPAddresses.DasIPAddresses, &dasIPAddress)
	}
	return &dasIPAddresses, nil
}

func (s *CommunicationServer) ChangeDASIPAddress(ctx context.Context, req *pb.DASIPAddress) (*pb.Ack, error) {
	var ack pb.Ack
	err := db.ChangeDASIP(req.Name, req.IpAddress)
	if err != nil {
		ack.Ok = false
		ack.Message = "Failed to change DAS IP address"
		return &ack, err
	}
	ack.Ok = true
	ack.Message = "DAS IP address changed successfully"
	return &ack, err
}

func (s *CommunicationServer) GetAcqRemarks(ctx context.Context, req *pb.ClientID) (*pb.AcqRemarks, error) {
	var acqRemarks pb.AcqRemarks
	acqRemarks.AcqRemarks = make([]*pb.AcqRemark, 0)
	acqDetail, err := db.GetLatestAcquisitionsList(true, 0)
	if err != nil {
		return nil, err
	}
	for _, acq := range acqDetail {
		var acqRemark pb.AcqRemark
		acqRemark.Date = acq.Date
		acqRemark.Time = acq.Time
		acqRemark.Remark = acq.Remark
		acqRemarks.AcqRemarks = append(acqRemarks.AcqRemarks, &acqRemark)
	}
	return &acqRemarks, nil
}

func (s *CommunicationServer) ChangeAcqRemark(ctx context.Context, req *pb.AcqRemark) (*pb.Ack, error) {
	var ack pb.Ack
	err := db.ChangeRemarks(req.Date, req.Time, req.Remark)
	if err != nil {
		ack.Ok = false
		ack.Message = "Failed to change acquisition remark"
		return &ack, err
	}
	ack.Ok = true
	ack.Message = "Acquisition remark changed successfully"
	return &ack, err
}
