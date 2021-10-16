rm -rf gochat

trap 'killall grpc-lab' SIGINT
go install -v

grpc-lab server --port=5001 --zkAddrs=localhost:2181 &
grpc-lab server --port=5002 --zkAddrs=localhost:2181 &
grpc-lab server --port=5003 --zkAddrs=localhost:2181 &

wait
