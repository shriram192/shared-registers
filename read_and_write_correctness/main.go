package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	//fmt.Print("Starting 50% read 50% write Workload.....")
	total_writes := 1000
	// total_keys := 1

	// Init Rand
	rand.Seed(time.Now().UnixNano())

	// batch_threshold := 10
	// start_throughput_timer := time.Now()

	for i := 1; i <= total_writes; i++ {
		fmt.Printf("Write Number: %d\n", i+1)
		get_random_write_key := strconv.Itoa(1)
		get_random_write_value := strconv.Itoa(rand.Intn(1000 + 1))

		//Exec Write Command
		write_args := []string{"set", get_random_write_key, get_random_write_value}
		write_cmd := exec.Command("./client", write_args...)
		write_abs_path, _ := filepath.Abs("../client")

		write_cmd.Dir = write_abs_path
		write_cmd.Path = "./client"

		write_cmd_output, write_cmd_err := write_cmd.CombinedOutput()
		if write_cmd_err != nil {
			fmt.Println(fmt.Sprint(write_cmd_err) + ": " + string(write_cmd_output))
		}

		fmt.Printf("Read Number: %d\n", i+1)
		get_random_read_key := strconv.Itoa(1)

		//Exec Read Command
		read_args := []string{"get", get_random_read_key}
		read_cmd := exec.Command("./client", read_args...)
		read_abs_path, _ := filepath.Abs("../client")

		read_cmd.Dir = read_abs_path
		read_cmd.Path = "./client"

		read_cmd_output, read_cmd_err := read_cmd.CombinedOutput()
		if read_cmd_err != nil {
			fmt.Println(fmt.Sprint(read_cmd_err) + ": " + string(read_cmd_output))
		}

	}
}
