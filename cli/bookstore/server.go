package bookstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct{
	port       int
	clientPort int

	mu    sync.RWMutex
	books []*bookv1.Book
	bookv1.UnimplementedBookServiceServer
}

func (s *Server) ListBooks(ctx context.Context, request *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &bookv1.ListBooksResponse{Books: s.books}, nil
}

func (s *Server) CreateBook(ctx context.Context, request *bookv1.CreateBookRequest) (*bookv1.CreateBookResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rand.Seed(time.Now().UnixMilli())
	book := &bookv1.Book{
		Id:     int64(rand.Intn(1000)),
		Author: request.Author,
		Title:  request.Title,
	}
	s.books = append(s.books, book)
	return &bookv1.CreateBookResponse{Book: book}, nil
}

func (s *Server) GetBook(ctx context.Context, request *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, book := range s.books {
		if book.Id == request.GetId() {
			return &bookv1.GetBookResponse{Book: book}, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *Server) DeleteBook(ctx context.Context, request *bookv1.DeleteBookRequest) (*bookv1.DeleteBookResponse, error) {
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
	return &bookv1.DeleteBookResponse{}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	bookv1.RegisterBookServiceServer(gs, s)
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
}

func Register(app *kingpin.Application) {
	c := new(Command)
	cmd := app.Command("bookstore", "Starts a bookstore server").Action(c.run)
	cmd.Flag("port", "Port on which the server should run").Required().IntVar(&c.port)
	cmd.Flag("clientPort", "This is where the gRPC-Gateway proxies the requests").Required().IntVar(&c.clientPort)
}

func (c *Command) run(*kingpin.ParseContext) error {
	s := &Server{port: c.port, clientPort: c.clientPort, books: make([]*bookv1.Book, 0)}
	return s.Start()
}
