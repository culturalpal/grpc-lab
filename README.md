# grpc-lab

This repository contains various experiment related to grpc in golang. The setup per experiment is described in the experiments section below 

The `contracts` directory contain a submodule which points to the [grpc-lab-contracts](https://github.com/hextechpal/grpc-lab-contracts) repo

The repo uses `buf` to generate the definitions. The contracts repo contains the `buf modules`. buf generate and work files are contained in this repo

You need to have the following `protobuf` tools in your path before you generate go structs from proto definitions

Protoc Gen for go

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

Protoc Gen for go-grpc

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Grpc Gateway tools

`go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest`

`go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest`


The generated files are contained in the generated directory. To regenerate the proto definitions use the following command
`buf generate`
## Experiments

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

### Http Transcoding using grpc gateway (Mongo go driver support as well)

Implemented a book service which is a basic in memory CRUD application. 
The primary thing is this book servce get exposes as an http service using `grpc-gateway` transcoding

To run this service you need to you need to build the repo using

`go build -o lab`

Once build you can run the bookstore example using 

`./lab bookstore --port=<port_of_grpc_server> --clientPort=<port_for_grpc_gateway_reverse_proxy>`

Once running you it exposes following http endpoints

- GET `/api/v1/books` - List of Books
- GET `/api/v1/books/{id}` - Get a book
- DELETE `/api/v1/books/{id}"` - Delete a book
- POST `/api/v1/books` - Create a book Payload - `'{"title":"Designing Data Intensive Applications", "author":"Michael Kleppman"}'`


