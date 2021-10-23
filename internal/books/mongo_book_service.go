package books

import (
	"context"
	"errors"
	bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"
	"github.com/ppal31/grpc-lab/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "books"

type MongoBookService struct {
	db *mongo.Database
	bookv1.UnimplementedBookServiceServer
}

func NewMongoBookService(db *mongo.Database) *MongoBookService {
	return &MongoBookService{db: db}
}

func (m *MongoBookService) ListBooks(ctx context.Context, request *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	cursor, err := m.db.Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var results []*bookv1.Book
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return &bookv1.ListBooksResponse{Books: results}, nil
}

func (m *MongoBookService) CreateBook(ctx context.Context, request *bookv1.CreateBookRequest) (*bookv1.CreateBookResponse, error) {
	book := &bookv1.Book{
		Id:     utils.GenerateUuid(),
		Author: &bookv1.Author{Name: request.Author},
		Title:  request.Title,
	}
	_, err := m.db.Collection(collectionName).InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}
	return &bookv1.CreateBookResponse{Book: book}, nil
}

func (m *MongoBookService) GetBook(ctx context.Context, request *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	sr := m.db.Collection(collectionName).FindOne(ctx, bson.M{"_id": request.Id})
	if sr.Err() != nil {
		return nil, sr.Err()
	}

	book := &bookv1.Book{}
	if err := sr.Decode(book); err != nil {
		return nil, err
	}
	return &bookv1.GetBookResponse{Book: book}, nil
}

func (m *MongoBookService) DeleteBook(ctx context.Context, request *bookv1.DeleteBookRequest) (*bookv1.DeleteBookResponse, error) {
	dr, err := m.db.Collection(collectionName).DeleteOne(ctx, bson.M{"_id": request.Id})
	if err != nil {
		return nil, err
	}

	if dr.DeletedCount < 1 {
		return nil, errors.New("book with id does not exist")
	}

	return &bookv1.DeleteBookResponse{}, nil
}
