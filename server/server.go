package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/shriram192/shared-registers/api"
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

	registers := make(map[string]string)

	for i := 1; i <= 100000; i++ {
		str_index := strconv.Itoa(i)
		registers[str_index] = "init"
	}

	s := api.Server{Registers: registers, Timestamp: 0}

	grpcServer := grpc.NewServer()

	api.RegisterApiServer(grpcServer, &s)

	fmt.Printf("Starting Server at %s\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
