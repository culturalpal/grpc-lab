package resolver

import (
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const (
	scheme   = "zk"
	rootPath = "/chat_servers_golang"
)

type ZkBuilder struct {
	zkAddrs []string
}

func NewZkBuilder(zkAddrs []string) *ZkBuilder {
	return &ZkBuilder{zkAddrs: zkAddrs}
}

func (zkb *ZkBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &ZkResolver{
		target:  target,
		cc:      cc,
		zkAddrs: zkb.zkAddrs,
	}
	err := r.start()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (zkb *ZkBuilder) Scheme() string {
	return scheme
}

type ZkResolver struct {
	target  resolver.Target
	cc      resolver.ClientConn
	zkAddrs []string
}

func (r *ZkResolver) start() error {
	zkc, _, err := zk.Connect(r.zkAddrs, 5*time.Second)
	if err != nil {
		return err
	}
	children, _, ch, err := zkc.ChildrenW(rootPath)
	srvs := buildAddresses(children)
	go r.updater(ch, zkc)
	return r.cc.UpdateState(resolver.State{Addresses: srvs})
}

func buildAddresses(children []string) []resolver.Address {
	srvs := make([]resolver.Address, len(children))
	for i, srv := range children {
		srvs[i] = resolver.Address{Addr: srv}
	}
	return srvs
}

func (r *ZkResolver) ResolveNow(o resolver.ResolveNowOptions) {

}

func (*ZkResolver) Close() {
	log.Printf("Close Called")
}

func (r *ZkResolver) updater(ch <-chan zk.Event, zkc *zk.Conn) {
	for {
		ev := <-ch
		if ev.Type == zk.EventNodeChildrenChanged {
			log.Printf("EventNodeChildrenChanged received updating servers")
			children, _, _ := zkc.Children(rootPath)
			srvs := buildAddresses(children)
			r.cc.UpdateState(resolver.State{Addresses: srvs})
		} else {
			//log.Printf("Ignoring event of the type %d", ev.Type)
		}

	}
}
