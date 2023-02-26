package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
	"os"
)

func main() {
	fmt.Print("Starting 50% read 50% write Workload.....")

	total_writes := 5000
	total_keys := 10000

	// Init Rand
	rand.Seed(time.Now().UnixNano())

	LOG_FILE := "../etcd_throughput_logs"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	batch_threshold := 10
	start_throughput_timer := time.Now()

	for i := 1; i <= total_writes; i++ {
		//fmt.Printf("Write Number: %d\n", i+1)

		get_random_write_key := strconv.Itoa(1 + rand.Intn(total_keys-1+1))
		get_random_write_value := strconv.Itoa(rand.Intn(1000 + 1))

		//Exec Write Command
		write_args := []string{"--endpoints=128.110.96.159:2380", "put", get_random_write_key, get_random_write_value}
		write_cmd := exec.Command("etcdctl", write_args...)
		write_cmd_output, write_cmd_err := write_cmd.CombinedOutput()
		if write_cmd_err != nil {
			fmt.Println(fmt.Sprint(write_cmd_err) + ": " + string(write_cmd_output))
		}

		//fmt.Printf("Read Number: %d\n", i+1)
		get_random_read_key := strconv.Itoa(1 + rand.Intn(total_keys-1+1))

		//Exec Read Command
		read_args := []string{"--endpoints=128.110.96.159:2380", "get", get_random_read_key}
		read_cmd := exec.Command("etcdctl", read_args...)

		read_cmd_output, read_cmd_err := read_cmd.CombinedOutput()
		if read_cmd_err != nil {
			fmt.Println(fmt.Sprint(read_cmd_err) + ": " + string(read_cmd_output))
		}

		if i % batch_threshold == 0 {
			end_throughput_timer := time.Now()
			elapsed := end_throughput_timer.Sub(start_throughput_timer)
			start_throughput_timer = time.Now()
			log.Printf("%f", float64(batch_threshold) / elapsed.Seconds())
		}
	}
}