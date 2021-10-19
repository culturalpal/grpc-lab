rm -rf grpc-lab

trap 'killall grpc-lab' SIGINT
go install -v

grpc-lab lbserver --port=5001 --zkAddrs=localhost:2181 &
grpc-lab lbserver --port=5002 --zkAddrs=localhost:2181 &
grpc-lab lbserver --port=5003 --zkAddrs=localhost:2181 &

wait
