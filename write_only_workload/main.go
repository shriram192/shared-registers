package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	log.Printf("Starting Write-Only Workload.....")
	start_time := time.Now()

	total_writes := 10000
	total_keys := 10000

	// Init Rand
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < total_writes; i++ {
		log.Printf("Write Number: %d", i+1)
		get_random_key := strconv.Itoa(1 + rand.Intn(total_keys-1+1))
		get_random_value := strconv.Itoa(rand.Intn(1000 + 1))

		//Exec Write Command

		args := []string{"set", get_random_key, get_random_value}
		cmd := exec.Command("./client", args...)
		abs_path, _ := filepath.Abs("../client")

		cmd.Dir = abs_path
		cmd.Path = "./client"

		cmd_output, cmd_err := cmd.CombinedOutput()
		if cmd_err != nil {
			fmt.Println(fmt.Sprint(cmd_err) + ": " + string(cmd_output))
		}
	}

	end_time := time.Now()
	elapsed := end_time.Sub(start_time)
	log.Printf("Total Time for Workload: %f", elapsed.Seconds())
}
