package main

import (
	"log"
	"os"

	"github.com/shriram192/shared-registers/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var servers = [5]string{"127.0.0.1:9000", "127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003", "127.0.0.1:9004"}

func getMaxTimestamp(a []int64) (max int64) {
	max = a[0]
	for _, value := range a {
		if value > max {
			max = value
		}
	}
	return max
}

func main() {

	if len(os.Args[1:]) < 2 {
		log.Fatalf("Incorrect call: go run write.go [key] [value] or ./writer [key] [value]")
	}

	var setKey = os.Args[1]
	var setVal = os.Args[2]

	var client_list []api.ApiClient

	var timestamp_list []int64

	for i := 0; i < len(servers); i++ {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(servers[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := api.NewApiClient(conn)
		client_list = append(client_list, c)

		read_res, read_err := c.GetValue(context.Background(), &api.ReadInput{Key: setKey})

		if read_err != nil {
			log.Fatalf("Error when calling GetValue: %v", err)
		}

		log.Printf("Read Value from Server %d: %s", i+1, read_res.Value)
		log.Printf("Read Timestamp from Server %d: %d", i+1, read_res.Timestamp)

		timestamp_list = append(timestamp_list, read_res.Timestamp)
	}

	var max_time_stamp = getMaxTimestamp(timestamp_list) + 1
	log.Printf("Max Time Stamp: %d ", max_time_stamp)

	for i := 0; i < len(client_list); i++ {

		write_res, write_err := client_list[i].PutValue(context.Background(), &api.WriteInput{Key: setKey, Value: setVal, Timestamp: max_time_stamp})
		if write_err != nil {
			log.Fatalf("Error when calling PutValue: %v", write_err)
		}

		log.Printf("Write Status from Server %d: %t", i+1, write_res.Status)
		log.Printf("Write Message %d: %s", i+1, write_res.Message)
	}
}
