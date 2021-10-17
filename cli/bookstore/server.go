package bookstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/ppal31/grpc-lab/generated/book"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	port       int
	clientPort int

	mu    sync.RWMutex
	books []*v1.Book
	v1.UnimplementedBookStoreServer
}

func (s *Server) ListBooks(ctx context.Context, request *v1.ListBookRequest) (*v1.ListBookResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &v1.ListBookResponse{Books: s.books}, nil
}

func (s *Server) CreateBook(ctx context.Context, request *v1.CreateBookRequest) (*v1.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rand.Seed(time.Now().UnixMilli())
	book := &v1.Book{
		Id:     int64(rand.Intn(1000)),
		Author: request.Author,
		Title:  request.Title,
	}
	s.books = append(s.books, book)
	return book, nil
}

func (s *Server) GetBook(ctx context.Context, request *v1.GetBookRequest) (*v1.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, book := range s.books {
		if book.Id == request.GetId() {
			return book, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *Server) DeleteBook(ctx context.Context, request *v1.DeleteBookRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delIdx := -1
	for idx, book := range s.books {
		if book.Id == request.GetId() {
			delIdx = idx
			break
		}
	}

	if delIdx < 0 {
		return nil, errors.New("book not found")
	}
	s.books = append(s.books[:delIdx], s.books[delIdx+1:]...)
	return nil, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	v1.RegisterBookStoreServer(gs, s)
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
	// Register Greeter
	err = v1.RegisterBookStoreHandler(context.Background(), gwmux, conn)
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
}

func Register(app *kingpin.Application) {
	c := new(Command)
	cmd := app.Command("bookstore", "Starts a bookstore server").Action(c.run)
	cmd.Flag("port", "Port on which the server should run").Required().IntVar(&c.port)
	cmd.Flag("clientPort", "This is where the gRPC-Gateway proxies the requests").Required().IntVar(&c.clientPort)
}

func (c *Command) run(*kingpin.ParseContext) error {
	s := &Server{port: c.port, clientPort: c.clientPort, books: make([]*v1.Book, 0)}
	return s.Start()
}
