build:
	go build -o client client/client.go
	go build -o replica replica/replica.go
	go build -o rw_workload read_and_write_workload/main.go
	go build -o r_workload read_only_workload/main.go
	go build -o w_workload write_only_workload/main.go

dep: 
	go mod download

run:
	./client

clean:
	go clean
