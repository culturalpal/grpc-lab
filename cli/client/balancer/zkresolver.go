package balancer

import (
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
	"time"
)

const (
	scheme   = "zk"
	rootPath = "/chat_servers_golang"
)

type ZkResolver struct {
	cc         resolver.ClientConn
	zkc        *zk.Conn
	serverList sync.Map
}

func NewZkBuilder(zkAddrs []string) (*ZkResolver, error) {
	zkc, _, err := zk.Connect(zkAddrs, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &ZkResolver{zkc: zkc}, nil
}

func (r *ZkResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	err := r.start()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ZkResolver) Scheme() string {
	return scheme
}

func (r *ZkResolver) start() error {
	children, _, ch, err := r.zkc.ChildrenW(rootPath)
	if err != nil {
		return err
	}
	r.updateServers(children)
	go r.updater(ch)
	return r.cc.UpdateState(resolver.State{Addresses: r.getServices()})
}

func (r *ZkResolver) updateServers(children []string) {
	for _, srv := range children {
		accountId, _, _ := r.zkc.Get(rootPath + "/" + srv)
		r.setServerList(srv, string(accountId))
	}
}

func (r *ZkResolver) setServerList(srv string, accountId string) {
	addr := resolver.Address{Addr: srv}
	addr = SetAddrInfo(addr, AddrInfo{accountId: accountId})
	r.serverList.Store(srv, addr)
	_ = r.cc.UpdateState(resolver.State{Addresses: r.getServices()})
	log.Println("put key :", srv, "accountId:", accountId)
}

func (r *ZkResolver) DelServiceList(key string) {
	r.serverList.Delete(key)
	_ = r.cc.UpdateState(resolver.State{Addresses: r.getServices()})
	log.Println("del key:", key)
}

func (r *ZkResolver) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	r.serverList.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}

func (r *ZkResolver) updater(ch <-chan zk.Event) {
	for {
		ev := <-ch
		if ev.Type == zk.EventNodeChildrenChanged {
			log.Printf("EventNodeChildrenChanged received updating servers")
			children, _, _ := r.zkc.Children(rootPath)
			r.updateServers(children)
		} else {
			//log.Printf("Ignoring event of the type %d", ev.Type)
		}

	}
}

func (r *ZkResolver) ResolveNow(o resolver.ResolveNowOptions) {

}

func (*ZkResolver) Close() {
	log.Printf("Close Called")
}
