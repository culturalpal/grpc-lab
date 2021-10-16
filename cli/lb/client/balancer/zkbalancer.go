package balancer

import (
	"errors"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

const Name = "account"

// attributeKey is the type used as the key to store AddrInfo in the Attributes
// field of resolver.Address.
type attributeKey struct{}

// AddrInfo will be stored inside Address metadata in order to use weighted balancer.
type AddrInfo struct {
	accountId string
}

// SetAddrInfo returns a copy of addr in which the Attributes field is updated
// with addrInfo.
func SetAddrInfo(addr resolver.Address, addrInfo AddrInfo) resolver.Address {
	addr.Attributes = attributes.New()
	addr.Attributes = addr.Attributes.WithValues(attributeKey{}, addrInfo)
	return addr
}

// GetAddrInfo returns the AddrInfo stored in the Attributes fields of addr.
func GetAddrInfo(addr resolver.Address) AddrInfo {
	v := addr.Attributes.Value(attributeKey{})
	ai, _ := v.(AddrInfo)
	return ai
}

type zkPickerBuilder struct{}

func (zpb *zkPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("zkPicker: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scsMap := make(map[string]balancer.SubConn)
	for subConn, addr := range info.ReadySCs {
		info := GetAddrInfo(addr.Address)
		scsMap[info.accountId] = subConn
	}
	return &zkPicker{
		subConns: scsMap,
	}
}

// NewBuilder creates a new weight balancer builder.
func NewBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &zkPickerBuilder{}, base.Config{HealthCheck: false})
}

type zkPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns map[string]balancer.SubConn

	mu sync.Mutex
}

func (p *zkPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	accountId := info.Ctx.Value("accountId").(string)
	for i := 0; i < 3; i++ {
		p.mu.Lock()
		sc, ok := p.subConns[accountId]
		p.mu.Unlock()
		if ok {
			return balancer.PickResult{SubConn: sc}, nil
		}
		time.Sleep(10 * time.Millisecond)
	}

	return balancer.PickResult{}, errors.New("server not available for account" + accountId)

}
