package api

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/sync/syncmap"
	"google.golang.org/grpc/status"
)

func latestState(s *api.Server) {
	LOG_FILE := "../log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error occured at Open: %v", err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	log_strings := make([]string, 0)

	for scanner.Scan() {
		log_strings = append(log_strings, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Read file Failure: %v", err)
	}

	if s.Timestamp < s.LogTimestamp {
		for _, str := range log_strings {
			line := strings.Split(str, " ")

			// Use Client ID and Sequence Number for Non-Commutative
			//clientId := line[4]
			//seqId := line[5]

			key := line[2]
			value := line[3]
			timestamp, _ := strconv.Atoi(line[6])

			// Store to KV Store
			if int64(timestamp) > s.Timestamp {
				fmt.Printf("Executing Write : Key: %s, Value: %s\n", key, value)
				s.Registers.Store(key, value)
				s.Timestamp = s.Timestamp + 1
			}
		}
		return
	} else {
		return
	}
}

// Server represents the gRPC server
type Server struct {
	Registers    syncmap.Map
	Timestamp    int64
	LogTimestamp int64
	UnimplementedApiServer
}

func (s *Server) GetState(ctx context.Context, in *StateInput) (*StateOutput, error) {
	map_data := make(map[string]string)
	s.Registers.Range(func(key any, value any) bool {
		map_data[key.(string)] = value.(string)
		return true
	})
	return &StateOutput{Registers: map_data, Timestamp: s.Timestamp, LogTimestamp: s.LogTimestamp}, nil
}

// Get Register Value and its timestamp Server RPC
func (s *Server) GetValue(ctx context.Context, in *ReadInput) (*ReadOutput, error) {
	log.Printf("Received Read Key: %s", in.Key)

	start_time := time.Now()
	// Bring the data store to latest config
	latestState(s)

	val, ok_val := s.Registers.Load(in.Key)
	if ok_val {
		end_time := time.Now()
		elapsed := end_time.Sub(start_time)
		log.Printf("R: %f", 1/elapsed.Seconds())

		return &ReadOutput{Value: val.(string)}, nil
	} else {
		end_time := time.Now()
		elapsed := end_time.Sub(start_time)
		log.Printf("R: %f", 1/elapsed.Seconds())
		return &ReadOutput{Value: "-1"}, status.Error(400, "Value not Found")
	}
}

// Write Register Value Server RPC
func (s *Server) PutValue(ctx context.Context, in *WriteInput) (*WriteOutput, error) {
	log.Printf("Received Write Key: %s", in.Key)
	log.Printf("Received Write Value: %s", in.Value)

	LOG_FILE := "../log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	s.LogTimestamp = s.LogTimestamp + 1
	log.Printf("%s %s %d %d %d", in.Key, in.Value, in.ClientId, in.SeqId, s.LogTimestamp)
	return &WriteOutput{Status: true, Message: "Item Stored in Log"}, nil
}
