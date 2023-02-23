package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/shriram192/shared-registers/api"
	"golang.org/x/sync/syncmap"
	"google.golang.org/grpc"
)

func main() {

	if len(os.Args[1:]) != 1 {
		log.Fatalf("Incorrect call: go run server.go [ip:port] or ./server [ip:port]")
	}

	var port = os.Args[1]

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var registers syncmap.Map
	var timestamps syncmap.Map

	total_keys := 10000

	for i := 1; i <= total_keys; i++ {
		str_index := strconv.Itoa(i)
		registers.Store(str_index, "init")
		timestamps.Store(str_index, int64(0))
	}

	s := api.Server{Registers: registers, Timestamps: timestamps}

	grpcServer := grpc.NewServer()

	api.RegisterApiServer(grpcServer, &s)

	fmt.Printf("Starting Server at %s\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
