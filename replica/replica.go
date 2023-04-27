package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shriram192/shared-registers/api"
	"golang.org/x/net/context"
	"golang.org/x/sync/syncmap"
	"google.golang.org/grpc"
)

func touchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func scpData(host string) {
	fmt.Println("Performing SCP to get log file!")
	fmt.Printf("SCP Host Name: %s\n", host)

	// scp sa84@10.10.1.3:~/go/src/github.com/shared-registers/log ../log

	op := "scp"
	user := "sa84"
	path := "~/go/src/github.com/shared-registers/log"

	source := user + "@" + host + ":" + path
	dest := "../log"

	cmd := exec.Command(op, source, dest)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("SCP Complete! Log Downloaded!")
	}

	return
}

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

func watchForLogs(filePath string, fn func(s *api.Server), s *api.Server) error {
	fmt.Printf("Initializing Watcher for Logging!")
	defer func() {
		r := recover()
		if r != nil {
			log.Fatalf("Error:watching file: %v", r)
		}
	}()
	initialStat, err := os.Stat(filePath)

	if err != nil {
		log.Fatalf("Error: Stat: %v", err)
	}

	for {
		stat, err := os.Stat(filePath)

		if err != nil {
			log.Fatalf("Error: Stat: %v", err)
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			fmt.Println("New Logs In Log File! Running Command!")
			fn(s)
			initialStat, err = os.Stat(filePath)

			if err != nil {
				log.Fatalf("Error: Stat: %v", err)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {

	if len(os.Args[1:]) > 1 {
		val, _ := strconv.Atoi(os.Args[1])
		if len(os.Args[3:]) < val {
			log.Fatalf("number of servers do not match number of ips provided")
		}
	} else {
		log.Fatalf("Incorrect call: go run server.go [num_servers] [idx] [**ip:port|replicas] or ./serve [num_servers] [idx] [**ip:port|replicas]")
	}

	idx, _ := strconv.Atoi(os.Args[2])
	replicas := os.Args[3:]
	port := replicas[idx]

	touchFile("../log")

	var maxTime int64 = -1
	var maxLogTime int64 = -1
	var maxId int = -1
	maxMap := make(map[string]string)

	for i := 0; i < len(replicas); i++ {
		if i != idx {

			var conn *grpc.ClientConn
			conn, err := grpc.Dial(replicas[i], grpc.WithInsecure())

			if err != nil {
				log.Printf("did not connect: %v %v", err, conn)
			}
			defer conn.Close()

			c := api.NewApiClient(conn)
			read_res, read_err := c.GetState(context.Background(), &api.StateInput{})

			if read_err != nil {
				log.Printf("Error when calling GetState: %v, for %s", err, replicas[i])
			} else {
				if read_res.Timestamp > maxTime {
					maxTime = read_res.Timestamp
					maxLogTime = read_res.LogTimestamp
					maxId = i
					maxMap = read_res.Registers
				}
			}

		}
	}

	var registers syncmap.Map
	var initTime = -1
	var initLogTime = -1
	if maxTime == -1 || maxTime == 0 {
		fmt.Println("Initializing Map! No Replica with Logs Found!")
		total_keys := 10000
		for i := 1; i <= total_keys; i++ {
			str_index := strconv.Itoa(i)
			registers.Store(str_index, "init")
		}
		initTime = 0
		initLogTime = 0
	} else {
		fmt.Println("Replicas found with Logs!")
		for k, v := range maxMap {
			registers.Store(k, v)
		}
		initTime = int(maxTime)
		initLogTime = int(maxLogTime)
		host_split := strings.Split(replicas[maxId], ":")
		scpData(host_split[0])
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.Server{Registers: registers, Timestamp: int64(initTime), LogTimestamp: int64(initLogTime)}

	grpcServer := grpc.NewServer()

	api.RegisterApiServer(grpcServer, &s)

	// Start Watcher Script
	go watchForLogs("../log", latestState, &s)

	fmt.Printf("Starting Server at %s\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
