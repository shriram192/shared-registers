package api

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server
type Server struct {
	Registers map[string]string
	Timestamp int64
	UnimplementedApiServer
}

// Get Register Value and its timestamp Server RPC
func (s *Server) GetValue(ctx context.Context, in *ReadInput) (*ReadOutput, error) {
	log.Printf("Received Read Key: %s", in.Key)

	val, ok := s.Registers[in.Key]
	if ok {
		return &ReadOutput{Value: val, Timestamp: s.Timestamp}, nil
	} else {
		return &ReadOutput{Value: "-1", Timestamp: -1}, status.Error(400, "Value not Found")
	}
}

// Write Register Value Server RPC
func (s *Server) PutValue(ctx context.Context, in *WriteInput) (*WriteOutput, error) {
	log.Printf("Received Write Key: %s", in.Key)
	log.Printf("Received Write Value: %s", in.Value)

	if in.Timestamp > s.Timestamp {
		s.Registers[in.Key] = in.Value
		s.Timestamp = in.Timestamp
		return &WriteOutput{Status: true, Message: "Item Stored in Register"}, nil
	}
	return &WriteOutput{Status: true, Message: "Timestamp provided is old"}, nil
}
