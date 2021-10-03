# gochat

This repository implemented a rpc server and client in go.

It also implements client side load balancing based on zookepper

### Setup

The repo contains a launcher script at the root ``start.sh``
The script assumes a running zookeeper instance on ``localhost:2181``. 

Its start 3 grpcs servers on port 5001, 5002 and 5003

The repo also contains a clinet which start with a zookeeper resolver in a round robin strategy

### Test

The repo contains an integration test ``ping_test.go`` which spins up a client which makes 100o pings to the server
In the logs you can see the requested getting routed to diff servers.

You can play around by killing and adding servers and the load balancer will take care of routing requests to accurate servers

### TODO

- [ ] Add better tests
- [ ] Add streaming demos to the repo
- [ ] Better error handling for zookeeper
- [ ] Add more examples using etcd/consul etc
- [ ] Better read me file


