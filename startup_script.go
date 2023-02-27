package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	log.Printf("Starting Write-Only Workload.....")
	start_time := time.Now()

	total_writes := 10000
	total_keys := 10000

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
