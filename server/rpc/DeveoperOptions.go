package rpc

import (
	"context"
	pb "csrsp/server/communication"
	"csrsp/server/global"
)

func (s *CommunicationServer) GetDeveloperOptions(ctx context.Context, req *pb.ClientID) (*pb.DeveloperOptions, error) {
	return &pb.DeveloperOptions{
		AutoArchival:             global.App.DeveloperOptions.AutomaticArchival,
		LogLevel:                 global.App.DeveloperOptions.LogLevel,
		EnableParallelProcessing: global.App.DeveloperOptions.ParallelProcessing,
		EncryptionMode:           global.App.DeveloperOptions.EncryptionMode,
		EndProcessID:             global.App.DeveloperOptions.EndProcessID,
		MaxThreads:               global.App.DeveloperOptions.MaxThreads,
		BufferLength:             global.App.DeveloperOptions.BufferLength,
	}, nil
}

func (s *CommunicationServer) SetDeveloperOptions(ctx context.Context, do *pb.DeveloperOptions) (*pb.Ack, error) {
	var dev global.DeveloperOptions
	dev.AutomaticArchival = do.AutoArchival
	dev.LogLevel = do.LogLevel
	dev.ParallelProcessing = do.EnableParallelProcessing
	dev.EncryptionMode = do.EncryptionMode
	dev.EndProcessID = do.EndProcessID
	dev.MaxThreads = do.MaxThreads
	dev.BufferLength = do.BufferLength
	global.SetDeveloperOptions(dev)
	err := dev.Save(global.App.Config.DevOpsPath)
	if err != nil {
		return &pb.Ack{
			Ok:      false,
			Message: "Unable to save Developer Options, these are applicable for this session only",
		}, nil
	}
	return &pb.Ack{
		Ok: true,
	}, nil
}
