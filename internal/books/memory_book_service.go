package books

import (
	"context"
	bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"
)

type MongoBookService struct {
	bookv1.UnimplementedBookServiceServer
}

func (m MongoBookService) ListBooks(ctx context.Context, request *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	panic("implement me")
}

func (m MongoBookService) CreateBook(ctx context.Context, request *bookv1.CreateBookRequest) (*bookv1.CreateBookResponse, error) {
	panic("implement me")
}

func (m MongoBookService) GetBook(ctx context.Context, request *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	panic("implement me")
}

func (m MongoBookService) DeleteBook(ctx context.Context, request *bookv1.DeleteBookRequest) (*bookv1.DeleteBookResponse, error) {
	panic("implement me")
}

func (m MongoBookService) mustEmbedUnimplementedBookServiceServer() {
	panic("implement me")
}


