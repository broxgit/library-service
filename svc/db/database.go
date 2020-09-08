package db

import "github.com/broxgit/library-service/internal/pkg/models"

type Database interface {
	// Create
	CreateBook(book models.Book) error

	// Read
	GetBook(bookId string) (models.Book, error)

	// Update
	UpdateBook(models.Book) error

	// Delete
	DeleteBook(bookId string) error

	// List
	ListBooks() ([]models.Book, error)
}
