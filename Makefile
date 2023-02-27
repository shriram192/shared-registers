build:
	go build -o client/client client/client.go
	go build -o replica/replica replica/replica.go
	go build -o read_and_write_workload/rw_workload read_and_write_workload/main.go
	go build -o read_only_workload/r_workload read_only_workload/main.go
	go build -o write_only_workload/w_workload write_only_workload/main.go

dep: 
	go mod download

run:
	./client

clean:
	go clean
