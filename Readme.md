This is the first version of our protocol Symi, without any ordering, this code assumes the workload is commutative and does not perform any ordering and performs execution directly. The algorithm mentioned in the report is under development and its implementation will be completed in the next few months.

How to run -

1. Cd replica | go build 
2. Cd client | go build
3. Run Replicas on each node server present with the command -  ./replica [num_servers] [current_ip] [ips (number same as num_servers)]
4. Run client with the command - ./client [set/get] [key] [value] [num_servers] [ips (number same as num_servers)]
5. Compile workload file - cd read_and_write_workload | go build 
6. Run Workload using command - ./read_and_write_workload [num_servers] [ips (number same as num_servers)]

Example commands - 
./replica 5 0 10.10.1.3:9000 10.10.1.1:9000 10.10.1.2:9000 10.10.1.4:9000 10.10.1.5:9000
./client set 1 1 5 10.10.1.3:9000 10.10.1.1:9000 10.10.1.2:9000 10.10.1.4:9000 10.10.1.5:9000
./read_and_write_workload 5 localhost:9000 localhost:9001 localhost:9002 localhost:9003 localhost:9004