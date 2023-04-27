package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shriram192/shared-registers/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func majorityElement(nums []string) string {
	lenNums := len(nums)

	if lenNums == 1 {
		return nums[0]
	}

	numsMap := make(map[string]int)

	for i := 0; i < lenNums; i++ {
		_, ok := numsMap[nums[i]]
		if ok {
			numsMap[nums[i]] = numsMap[nums[i]] + 1
			if numsMap[nums[i]] > lenNums/2 {
				return nums[i]
			}
		} else {
			numsMap[nums[i]] = 1
		}
	}

	return "-1"

}

func main() {

	if len(os.Args[1:]) > 1 {
		if os.Args[1] == "get" {
			val, _ := strconv.Atoi(os.Args[3])
			if len(os.Args[4:]) < val {
				log.Fatalf("number of servers do not match number of ips provided")
			}
		} else {
			val, _ := strconv.Atoi(os.Args[4])
			if len(os.Args[5:]) < val {
				log.Fatalf("number of servers do not match number of ips provided")
			}
		}
	} else {
		log.Fatalf("Incorrect call: ./client [get] [key] [num_servers] [**ip:port] or ./client [set] [key] [value] [num_servers] [**ip:port]")
	}

	operation := os.Args[1]
	var num_servers = 0
	var replicas = make([]string, 0)

	if operation == "get" {
		num_servers, _ = strconv.Atoi(os.Args[3])
		replicas = os.Args[4:]
	} else {
		num_servers, _ = strconv.Atoi(os.Args[4])
		replicas = os.Args[5:]
	}

	// LOG_FILE := "../client_logs"
	// logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)
	// log.SetFlags(log.LstdFlags)

	if operation == "get" {
		getKey := os.Args[2]

		var val_list []string

		start_time := time.Now()

		for i := 0; i < num_servers; i++ {
			var conn *grpc.ClientConn
			conn, err := grpc.Dial(replicas[i], grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()

			c := api.NewApiClient(conn)

			read_res, read_err := c.GetValue(context.Background(), &api.ReadInput{Key: getKey})

			if read_err != nil {
				log.Fatalf("Error when calling GetValue: %v", err)
			} else {
				val_list = append(val_list, read_res.Value)
			}
		}

		major_val := majorityElement(val_list)
		if major_val != "-1" {
			fmt.Printf("Majority Found!!! Read Value: %s\n", major_val)
		} else {
			fmt.Println("Majority Not Found!! Failure!!")
		}

		end_time := time.Now()
		latency := end_time.Sub(start_time)
		log.Printf("R: %f", latency.Seconds())

	} else if operation == "set" {
		setKey := os.Args[2]
		setVal := os.Args[3]

		start_time := time.Now()

		for i := 0; i < num_servers; i++ {
			var conn *grpc.ClientConn
			conn, err := grpc.Dial(replicas[i], grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()

			c := api.NewApiClient(conn)

			write_res, write_err := c.PutValue(context.Background(), &api.WriteInput{Key: setKey, Value: setVal})
			if write_err != nil {
				log.Fatalf("Error when calling PutValue: %v", write_err)
			} else {
				fmt.Printf("Status: %s\n", write_res.Message)
			}
		}

		end_time := time.Now()
		latency := end_time.Sub(start_time)

		log.Printf("W: %f", latency.Seconds())

	} else {
		log.Fatalf("Invalid Operation %s", os.Args[1])
	}
}
