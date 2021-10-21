package books

import (
	"context"
	"errors"
	bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"
	"github.com/ppal31/grpc-lab/internal/utils"
	"sync"
)

type MemoryBookService struct {
	mu    sync.RWMutex
	books []*bookv1.Book
	bookv1.UnimplementedBookServiceServer
}

func NewMemoryBookService(books []*bookv1.Book) *MemoryBookService {
	return &MemoryBookService{books: books}
}

func (mem *MemoryBookService) ListBooks(ctx context.Context, request *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	return &bookv1.ListBooksResponse{Books: mem.books}, nil
}

func (mem *MemoryBookService) CreateBook(ctx context.Context, request *bookv1.CreateBookRequest) (*bookv1.CreateBookResponse, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()
	book := &bookv1.Book{
		Id:     utils.GenerateUuid(),
		Author: request.Author,
		Title:  request.Title,
	}
	mem.books = append(mem.books, book)
	return &bookv1.CreateBookResponse{Book: book}, nil
}

func (mem *MemoryBookService) GetBook(ctx context.Context, request *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	for _, book := range mem.books {
		if book.Id == request.GetId() {
			return &bookv1.GetBookResponse{Book: book}, nil
		}
	}
	return nil, errors.New("book not found")
}

func (mem *MemoryBookService) DeleteBook(ctx context.Context, request *bookv1.DeleteBookRequest) (*bookv1.DeleteBookResponse, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()
	delIdx := -1
	for idx, book := range mem.books {
		if book.Id == request.GetId() {
			delIdx = idx
			break
		}
	}

	if delIdx < 0 {
		return nil, errors.New("book not found")
	}
	mem.books = append(mem.books[:delIdx], mem.books[delIdx+1:]...)
	return &bookv1.DeleteBookResponse{}, nil
}
