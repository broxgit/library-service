package db

import (
	"errors"
	"github.com/broxgit/library-service/internal/pkg/models"
)

// CacheDatabase implements the Database interface
// the provided database is a map that lives in memory and dies with the service/container
type CacheDatabase struct {
	bookCache map[string]models.Book
}

// create a new CacheDatabase 'object'
func newCacheDatabase() (*CacheDatabase, error) {
	db := &CacheDatabase{bookCache: make(map[string]models.Book)}
	return db, nil
}

// add new book to cache
func (db *CacheDatabase) CreateBook(book models.Book) error {
	db.bookCache[book.Id] = book
	return nil
}

// get book from cache
func (db *CacheDatabase) GetBook(bookId string) (models.Book, error) {
	book, found := db.bookCache[bookId]
	if !found {
		return book, errors.New("book was not found in the database")
	}
	return book, nil
}

// update book in cache
func (db *CacheDatabase) UpdateBook(book models.Book) error {
	db.bookCache[book.Id] = book
	return nil
}

// delete book from cache
func (db *CacheDatabase) DeleteBook(bookId string) error {
	delete(db.bookCache, bookId)
	return nil
}

// get all books from cache
func (db *CacheDatabase) ListBooks() ([]models.Book, error) {
	var books []models.Book
	for _, book := range db.bookCache {
		books = append(books, book)
	}
	return books, nil
}
