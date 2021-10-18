package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/ppal31/grpc-lab/cli/lb"
	pb "github.com/ppal31/grpc-lab/generated/chat/v1"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Server struct {
	zkAddrs []string
	port    int
	pb.UnimplementedChatServiceServer
}

func (s *Server) Start() error {
	if err := registerWithZk(s.zkAddrs, "localhost", s.port); err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterChatServiceServer(gs, s)
	log.Printf("server listening at %v", lis.Addr())
	return gs.Serve(lis)
}

func registerWithZk(zkAddrs []string, serverIp string, port int) error {
	zkc, _, err := zk.Connect(zkAddrs, 5*time.Second)
	if err != nil {
		return err
	}

	exists, _, err := zkc.Exists(lb.RootPath)
	if err != nil {
		return err
	}

	if !exists {
		_, err := zkc.Create(lb.RootPath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	//Root path exists or created hence update the child
	sp := lb.RootPath + "/" + serverIp + ":" + strconv.Itoa(port)
	_, err = zkc.Create(sp, []byte(strconv.Itoa(port)), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Chat(stream pb.ChatService_ChatServer) error {
	for {
		cm, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Message Received %s\n", cm.Message)
		stream.Send(&pb.ChatMessage{Message: fmt.Sprintf("%s Acked from %d", cm.Message, s.port)})
	}
	return nil
}

func (s *Server) Ping(ctx context.Context, pr *pb.PingRequest) (*pb.PongReply, error) {
	return &pb.PongReply{Message: fmt.Sprintf("Pong from %v", s.port)}, nil
}

type Command struct {
	port    int
	zkAddrs []string
}

func Register(app *kingpin.Application) {
	c := new(Command)
	cmd := app.Command("lbserver", "Starts a chat server").Action(c.run)
	cmd.Flag("port", "Port on which the server should run").Required().IntVar(&c.port)
	cmd.Flag("zkAddrs", "Address of zookeeper server").Required().StringsVar(&c.zkAddrs)
}

func (c *Command) run(*kingpin.ParseContext) error {
	s := &Server{port: c.port, zkAddrs: c.zkAddrs}
	return s.Start()
}
