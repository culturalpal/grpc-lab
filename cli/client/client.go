package client

import (
	"context"
	"fmt"
	pb "github.com/ppal31/gochat/api"
	gb "github.com/ppal31/gochat/cli/client/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
	"io"
	"log"
	"time"
)

type Client struct {
	zkAddrs []string // Address of server to connect to
}

func (c *Client) Ping() error {
	cc, closer, err := c.setupClient()
	if err != nil {
		return err
	}
	defer closer.Close()

	for i := 0; i < 1000; i++ {
		accountId := "5001"
		if i%2 == 0 {
			accountId = "5002"
		}
		ctx := context.WithValue(context.Background(), "accountId", accountId)
		r, err := cc.Ping(ctx, &pb.PingRequest{Message: "PING"})
		if err != nil {
			return err
		}
		log.Printf("Reply : %s", r.GetMessage())
		time.Sleep(time.Second)
	}
	return nil
}

func (c *Client) Chat(accountId string) error {
	cc, closer, err := c.setupClient()
	if err != nil {
		return err
	}
	defer closer.Close()
	ctx := context.WithValue(context.Background(), "accountId", accountId)
	cs, err := cc.Chat(ctx)
	if err != nil {
		return err
	}
	done := make(chan struct{})

	go func() {
		for {
			cm, err := cs.Recv()
			if err == io.EOF {
				done <- struct{}{}
				return
			}
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf(cm.Message)
		}

	}()

	for i := 0; i < 100; i++ {
		err := cs.Send(&pb.ChatMessage{Message: fmt.Sprintf("Message %d", i)})
		if err != nil {
			return err
		}
	}
	cs.CloseSend()
	<-done
	return nil
}

func (c *Client) setupClient() (pb.ChatServiceClient, io.Closer, error) {
	balancer.Register(gb.NewBuilder())
	rb, err := gb.NewZkBuilder(c.zkAddrs)
	if err != nil {
		return nil, nil, err
	}
	resolver.Register(rb)
	conn, err := grpc.Dial(fmt.Sprintf("%s:///", rb.Scheme()), grpc.WithInsecure(), grpc.WithBalancerName(gb.Name))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cc := pb.NewChatServiceClient(conn)
	return cc, conn, nil
}

func NewClient(zkAddrs []string) *Client {
	return &Client{zkAddrs: zkAddrs}
}
