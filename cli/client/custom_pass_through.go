package client

import (
	"google.golang.org/grpc/resolver"
	"log"
	"math/rand"
	"time"
)

const scheme = "customPassThrough"

type customPassThroughBuilder struct{}

func (*customPassThroughBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &CustomPassThroughResolver{
		target: target,
		cc:     cc,
	}
	r.start()
	return r, nil
}

func (*customPassThroughBuilder) Scheme() string {
	return scheme
}

type CustomPassThroughResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

func (r *CustomPassThroughResolver) start() {
	r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
}

func (r *CustomPassThroughResolver) ResolveNow(o resolver.ResolveNowOptions) {
	rand.Seed(time.Now().UnixNano())
	adds := []string{"localhost:50051", "localhost:50052"}
	add := adds[rand.Intn(2)]
	log.Printf("Address Selected %v\n", add)
	r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: add}}})
}

func (*CustomPassThroughResolver) Close() {
	log.Printf("Close Called")
}
