package client

import (
	"context"
	pb "github.com/ppal31/gochat/api"
	"github.com/ppal31/gochat/cli/client/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"time"
)

type Client struct {
	zkAddrs []string // Address of server to connect to
}

func (c *Client) Ping() error {
	conn, err := grpc.Dial("zk:///", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithResolvers(resolver.NewZkBuilder(c.zkAddrs)), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cc := pb.NewChatServiceClient(conn)

	for i := 0; i < 1000; i++ {
		r, err := cc.Ping(context.Background(), &pb.PingRequest{Message: "PING"})
		if err != nil {
			return err
		}
		log.Printf("Reply : %s", r.GetMessage())
		time.Sleep(time.Second)
	}
	return nil
}

func NewClient(zkAddrs []string) *Client {
	return &Client{zkAddrs: zkAddrs}
}
