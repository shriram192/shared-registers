package api

import (
	"fmt"
	"log"
	"os"
	"time"
	"encoding/json"
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
	//log.Printf("Received Read Key: %s", in.Key)

	LOG_FILE := "../throughput_logs"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	start_time := time.Now()
	val, ok_val := s.Registers.Load(in.Key)
	timestamp, ok_timestamp := s.Timestamps.Load(in.Key)
	if ok_val && ok_timestamp {
		end_time := time.Now()
		elapsed := end_time.Sub(start_time)
		bs, _ := json.Marshal(Registers)
    	log.Println(string(bs))
		// log.Printf("R: %f", 1/elapsed.Seconds())
		//fmt.Printf("R: %f\n", 1/elapsed.Seconds())

		return &ReadOutput{Value: val.(string), Timestamp: timestamp.(int64)}, nil
	} else {
		end_time := time.Now()
		elapsed := end_time.Sub(start_time)
		bs, _ := json.Marshal(Registers)
    	log.Println(string(bs))
		// log.Printf("R: %f", 1/elapsed.Seconds())
		//fmt.Printf("R: %f\n", 1/elapsed.Seconds())
		return &ReadOutput{Value: "-1", Timestamp: -1}, status.Error(400, "Value not Found")
	}
}

// Write Register Value Server RPC
func (s *Server) PutValue(ctx context.Context, in *WriteInput) (*WriteOutput, error) {
	//log.Printf("Received Write Key: %s", in.Key)
	//log.Printf("Received Write Value: %s", in.Value)

	LOG_FILE := "../throughput_logs"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	start_time := time.Now()
	timestamp, ok_timestamp := s.Timestamps.Load(in.Key)

	if ok_timestamp {
		if in.Timestamp > timestamp.(int64) {
			s.Registers.Store(in.Key, in.Value)
			s.Timestamps.Store(in.Key, in.Timestamp)
			end_time := time.Now()
			elapsed := end_time.Sub(start_time)
			bs, _ := json.Marshal(Registers)
    		log.Println(string(bs))
			// log.Printf("W: %f", 1/elapsed.Seconds())
			// fmt.Printf("W: %f\n", 1/elapsed.Seconds())
			return &WriteOutput{Status: true, Message: "Item Stored in Register"}, nil
		} else {
			end_time := time.Now()
			elapsed := end_time.Sub(start_time)
			bs, _ := json.Marshal(Registers)
    		log.Println(string(bs))
			// log.Printf("W: %f", 1/elapsed.Seconds())
			// fmt.Printf("W: %f\n", 1/elapsed.Seconds())
			return &WriteOutput{Status: true, Message: "Timestamp is old"}, nil
		}

	} else {
		end_time := time.Now()
		elapsed := end_time.Sub(start_time)
		bs, _ := json.Marshal(Registers)
    	log.Println(string(bs))
		// log.Printf("W: %f", 1/elapsed.Seconds())
		// fmt.Printf("W: %f\n", 1/elapsed.Seconds())
		return &WriteOutput{Status: false, Message: "Key not found"}, nil
	}

}
