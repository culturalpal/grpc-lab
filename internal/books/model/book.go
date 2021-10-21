package model

import bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"

type Book struct {
	Id     string `json:"id" bson:"_id"`
	Author string `json:"author" bson:"author"`
	Title  string `json:"title" bson:"title"`
}

func ToBook(book *Book) *bookv1.Book {
	return &bookv1.Book{
		Id:     book.Id,
		Author: book.Author,
		Title:  book.Title,
	}
}
