package rpc

import (
	"fmt"
	"log"
	"time"

	pb "csrsp/server/communication"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type CommunicationServer struct {
	pb.UnimplementedCommunicationServer
}

func (s *CommunicationServer) GetServerStatus(req *pb.ClientID, stream pb.Communication_GetServerStatusServer) error {
	log.Printf("Starting status stream for client: %s", req.Id)
	// Stream data every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Client %s disconnected", req.Id)
			return nil
		case t := <-ticker.C:
			// Get actual system stats
			v, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("Error getting memory stats: %v", err)
				continue
			}

			c, err := cpu.Percent(0, false)
			if err != nil {
				log.Printf("Error getting cpu stats: %v", err)
				continue
			}

			// cpu.Percent returns a slice, we want the total (first element)
			cpuUsage := 0.0
			if len(c) > 0 {
				cpuUsage = c[0]
			}

			status := &pb.ServerStatus{
				Timestamp: t.Format(time.RFC3339),
				Memory:    v.UsedPercent,
				Cpu:       cpuUsage,
			}
			if err := stream.Send(status); err != nil {
				return fmt.Errorf("failed to send status: %v", err)
			}
		}
	}
}
