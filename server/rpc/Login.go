package rpc

import (
	"context"
	"log"
	"net"
	"slices"

	pb "csrsp/server/communication"
	"csrsp/server/global"

	"google.golang.org/grpc/peer"
)

// IsWhitelisted checks if the client's IP address is allowed to access the server.
func (s *CommunicationServer) IsWhitelisted(ctx context.Context, req *pb.ClientID) (*pb.IsWhitelistedResponse, error) {
	var clientIP string
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
		// Attempt to strip port if present to get just the IP
		host, _, err := net.SplitHostPort(clientIP)
		if err == nil {
			clientIP = host
		}
	}

	log.Printf("Received whitelist check from IP: %s for ClientID: %s", clientIP, req.Id)

	isAllowed := slices.Contains(global.App.WhiteListedIPs, clientIP)

	if !isAllowed {
		log.Printf("Access denied for IP: %s", clientIP)
	} else {
		log.Printf("Access granted for IP: %s", clientIP)
	}

	return &pb.IsWhitelistedResponse{
		Whitelisted: isAllowed,
	}, nil
}
