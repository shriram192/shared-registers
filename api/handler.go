package api

import (
	"log"

	"golang.org/x/net/context"
	"golang.org/x/sync/syncmap"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server
type Server struct {
	Registers  syncmap.Map
	Timestamps syncmap.Map
	UnimplementedApiServer
}

// Get Register Value and its timestamp Server RPC
func (s *Server) GetValue(ctx context.Context, in *ReadInput) (*ReadOutput, error) {
	log.Printf("Received Read Key: %s", in.Key)

	val, ok_val := s.Registers.Load(in.Key)
	timestamp, ok_timestamp := s.Timestamps.Load(in.Key)
	if ok_val && ok_timestamp {
		return &ReadOutput{Value: val.(string), Timestamp: timestamp.(int64)}, nil
	} else {
		return &ReadOutput{Value: "-1", Timestamp: -1}, status.Error(400, "Value not Found")
	}
}

// Write Register Value Server RPC
func (s *Server) PutValue(ctx context.Context, in *WriteInput) (*WriteOutput, error) {
	log.Printf("Received Write Key: %s", in.Key)
	log.Printf("Received Write Value: %s", in.Value)

	timestamp, ok_timestamp := s.Timestamps.Load(in.Key)

	if ok_timestamp {
		if in.Timestamp > timestamp.(int64) {
			s.Registers.Store(in.Key, in.Value)
			s.Timestamps.Store(in.Key, in.Timestamp)
			return &WriteOutput{Status: true, Message: "Item Stored in Register"}, nil
		} else {
			return &WriteOutput{Status: true, Message: "Timestamp is old"}, nil
		}

	} else {
		return &WriteOutput{Status: false, Message: "Key not found"}, nil
	}

}
