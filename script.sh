num_replicas=$1
num_client=$2

start_time="$(date +%s)"
for (( i=1 ; i<=$num_replicas ; i++ )); 
do
    ./replica/replica 127.0.0.1:$((8999+i)) &
done

sleep 5

ops=`echo "scale=0 ; 10240 / $num_client" | bc`

cd read_and_write_workload/

for (( i=1 ; i<=$num_client ; i++ )); 
do
    ./rw_workload &
done

sleep 10

client_pid=`pgrep rw_workload`

for pid in $client_pid
do
    kill -9 $pid
done

server_pid=`pgrep replica`

for pid in $server_pid
do
    kill -9 $pid
done
end_time="$(date +%s)"

elapsed="$(bc <<<"$end_time-$start_time")"
echo "Total of $elapsed seconds elapsed for process"

