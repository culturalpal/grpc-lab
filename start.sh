rm -rf gochat

trap 'killall gochat' SIGINT
go install -v

gochat server --port=5001 --zkAddrs=localhost:2181 &
gochat server --port=5002 --zkAddrs=localhost:2181 &
gochat server --port=5003 --zkAddrs=localhost:2181 &

wait
