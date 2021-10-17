# grpc-lab

This repository contains various experiment related to grpc in golang. The setup per experiment is described in the setup section below 

The `contracts` directory contain a submodule which points to the [grpc-lab-contracts](https://www.google.com) repo

The generated files are contained in the generated directory. To regenerate the proto definitions use the following command

`protoc --go_out=./generated --go_opt=paths=source_relative  --go-grpc_out=./generated --go-grpc_opt=paths=source_relative <path/to/proto/in/contracts>`
## Setup

### Client Side Load Balancing experiment
The repo contains a launcher script at the root ``start.sh``
The script assumes a running zookeeper instance on ``localhost:2181``. 

Its start 3 grpcs servers on port 5001, 5002 and 5003

The repo also contains a clinet which start with a zookeeper resolver in a round robin strategy

#### Test

The repo contains an integration test ``balancer_test.go`` 
It contain two tests
- PingTest : This test sping up a clinet and make 100 pings to the server based on different `accountId` (just some metadata). Based on this metadata the balancer chooses which server the ping should be routed too
- ChatTest : It works same as ping test except the rpc used underneath uses bi-di streaming

#### TODO
- [ ] Add better tests
- [ ] Better error handling for zookeeper
- [ ] Add more examples using etcd/consul etc

### Implement a chat server


