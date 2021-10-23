package bookstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"
	"github.com/ppal31/grpc-lab/internal/books"
	"github.com/ppal31/grpc-lab/internal/storage/mongo"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net"
	"net/http"
)

type Server struct {
	port       int
	clientPort int

	bs bookv1.BookServiceServer
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	bookv1.RegisterBookServiceServer(gs, s.bs)
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		err := gs.Serve(lis)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", s.port),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = bookv1.RegisterBookServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.clientPort),
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://0.0.0.0:%d\n", s.clientPort)
	return gwServer.ListenAndServe()
}

type Command struct {
	port       int
	clientPort int
	backend    string
}

func Register(app *kingpin.Application) {
	c := new(Command)
	cmd := app.Command("bookstore", "Starts a bookstore server").Action(c.run)
	cmd.Flag("port", "Port on which the server should run").Required().IntVar(&c.port)
	cmd.Flag("clientPort", "This is where the gRPC-Gateway proxies the requests").Required().IntVar(&c.clientPort)
	cmd.Flag("backend", "The storage backend to use. The options currently are inmemory and mongo").Required().StringVar(&c.backend)
}

func (c *Command) run(*kingpin.ParseContext) error {
	switch c.backend {
	case "inmemory":
		s := &Server{port: c.port, clientPort: c.clientPort, bs: books.NewMemoryBookService([]*bookv1.Book{})}
		return s.Start()
	case "mongo":
		db, err := mongo.InitDb("mongodb://localhost:27017", "grpc-lab")
		if err != nil {
			return err
		}
		s := &Server{port: c.port, clientPort: c.clientPort, bs: books.NewMongoBookService(db)}
		return s.Start()
	default:
		return errors.New("please provide correct backend" + c.backend)
	}

}
