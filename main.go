package main

import (
	"log"
	"net"
	"sync"

	"github.com/go-training/grpc-health-check/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	port = "9000"
)

// Server is used to implement gorush grpc server.
type Server struct {
	mu sync.Mutex
	// statusMap stores the serving status of the services this Server monitors.
	statusMap map[string]proto.HealthCheckResponse_ServingStatus
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		statusMap: make(map[string]proto.HealthCheckResponse_ServingStatus),
	}
}

// Check implements `service Health`.
func (s *Server) Check(ctx context.Context, in *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if in.Service == "" {
		// check the server overall health status.
		return &proto.HealthCheckResponse{
			Status: proto.HealthCheckResponse_SERVING,
		}, nil
	}
	if status, ok := s.statusMap[in.Service]; ok {
		return &proto.HealthCheckResponse{
			Status: status,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}

func main() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := NewServer()
	proto.RegisterHealthServer(s, srv)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("gRPC server is running on " + port + " port.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
