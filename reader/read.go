package main

import (
	"log"
	"os"

	"github.com/shriram192/shared-registers/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var servers = [5]string{"127.0.0.1:9000", "127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003", "127.0.0.1:9004"}

func getMaxTimestampAndVal(a []int64, b []string) (max_timestamp int64, max_val string) {
	max_timestamp = a[0]
	max_val = b[0]
	for id, value := range a {
		if value > max_timestamp {
			max_timestamp = value
			max_val = b[id]
		}
	}
	return max_timestamp, max_val
}

func main() {

	if len(os.Args[1:]) == 0 {
		log.Fatalf("Incorrect call: go run read.go [key] [value] or ./reader [key] [value]")
	}

	var getKey = os.Args[1]

	var client_list []api.ApiClient

	var timestamp_list []int64
	var val_list []string

	for i := 0; i < len(servers); i++ {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(servers[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := api.NewApiClient(conn)
		client_list = append(client_list, c)

		read_res, read_err := c.GetValue(context.Background(), &api.ReadInput{Key: getKey})

		if read_err != nil {
			log.Fatalf("Error when calling GetValue: %v", err)
		}

		log.Printf("Read Value from Server %d: %s", i+1, read_res.Value)
		log.Printf("Read Timestamp from Server %d: %d", i+1, read_res.Timestamp)

		timestamp_list = append(timestamp_list, read_res.Timestamp)
		val_list = append(val_list, read_res.Value)
	}

	var max_time_stamp, max_val = getMaxTimestampAndVal(timestamp_list, val_list)
	log.Printf("Max Time Stamp: %d", max_time_stamp)
	log.Printf("Max Time Stamp Val: %s", max_val)

	for i := 0; i < len(client_list); i++ {

		write_res, write_err := client_list[i].PutValue(context.Background(), &api.WriteInput{Key: getKey, Value: max_val, Timestamp: max_time_stamp})
		if write_err != nil {
			log.Fatalf("Error when calling PutValue: %v", write_err)
		}

		log.Printf("Write Status from Server %d: %t", i+1, write_res.Status)
		log.Printf("Write Message %d: %s", i+1, write_res.Message)
	}
}
