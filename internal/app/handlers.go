package app

import (
	"github.com/broxgit/library-service/internal/pkg/errors"
	"github.com/broxgit/library-service/internal/pkg/models"
	"github.com/broxgit/library-service/svc/db"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

// initialize with a database object
func init() {
	var err error

	bookDatabase, err = db.GetDb()

	if err != nil {
		log.Fatal(err.Error())
	}
}

var bookDatabase db.Database

// CreateBook - Create a book
func CreateBook(c *gin.Context) {
	book := models.Book{}

	err := c.BindJSON(&book)
	if err != nil {
		handleError(c, err, errors.JSON_PARSE_ERROR())
		return
	}

	// check if book already exists
	books, err := bookDatabase.ListBooks()
	for _, _book := range books {
		if compareBooks(book, _book) {
			handleError(c, nil, errors.BOOK_ALREADY_EXISTS(_book.Id))
			return
		}
	}

	// create book metadata
	book.CreationTime = time.Now()
	book.LastUpdateTime = time.Now()
	book.Id = guuid.New().String()
	book.Version = guuid.New().String()
	err = bookDatabase.CreateBook(book)
	if err != nil {
		handleError(c, err, errors.BOOK_SAVE_ERROR())
	}
	c.JSON(http.StatusCreated, book)
}

// DeleteBookById - Deletes a book for a given id.
func DeleteBookById(c *gin.Context) {
	// delete is a no-op if the element doesn't exist,
	// so we shouldn't receive an error

	book, err := getBookFromDatabase(c)
	if err == nil {
		// verify that If-Match equals current version
		if book.Version == c.GetHeader("If-Match") {
			err := bookDatabase.DeleteBook(c.Param("id"))
			if err != nil {
				handleError(c, err, errors.INTERNAL_SERVER_ERROR())
			}
			c.JSON(http.StatusNoContent, gin.H{})
		} else {
			handleError(c, nil, errors.INVALID_IF_MATCH(book.Version, c.GetHeader("If-Match")))
			return
		}
	} else {
		// book didn't exist in the cache, return an error
		handleError(c, err, errors.BOOK_NOT_FOUND(c.Param("id")))
		return
	}
}

// GetBookById - Returns a book for a given id.
func GetBookById(c *gin.Context) {
	book, err := getBookFromDatabase(c)
	if err == nil {
		c.JSON(http.StatusOK, book)
	} else {
		// book didn't exist in the cache, return an error
		handleError(c, err, errors.BOOK_NOT_FOUND(c.Param("id")))
		return
	}
}

// ListBooks - List all books
func ListBooks(c *gin.Context) {
	var books []models.Book

	// initialize empty array in the case that we don't have any cached books
	books = make([]models.Book, 0)

	// build the 'books' array
	books, err := bookDatabase.ListBooks()

	if err != nil {
		handleError(c, err, errors.INTERNAL_SERVER_ERROR())
	}

	c.JSON(http.StatusOK, books)
}

// UpdateBookById - Updates a book for a given id
func UpdateBookById(c *gin.Context) {
	book := models.Book{}

	err := c.BindJSON(&book)
	if err != nil {
		handleError(c, err, errors.JSON_PARSE_ERROR())
		return
	}

	bookCheck, err := getBookFromDatabase(c)

	// if book doesn't exist in cache, return an error
	if err != nil {
		handleError(c, nil, errors.BOOK_NOT_FOUND(c.Param("id")))
		return
	}

	// verify that If-Match equals current version
	ifMatch := c.GetHeader("If-Match")
	if bookCheck.Version == ifMatch {
		// update book metadata
		book.LastUpdateTime = time.Now()
		book.Version = guuid.New().String()
		err := bookDatabase.UpdateBook(book)
		if err != nil {
			// an unexpected error has occurred
			handleError(c, err, errors.INTERNAL_SERVER_ERROR())
		}

		c.JSON(http.StatusOK, book)
	} else {
		handleError(c, nil, errors.INVALID_IF_MATCH(book.Version, ifMatch))
		return
	}
}

// generic error handling function
func handleError(c *gin.Context, err error, libError errors.LibraryError) {
	log.Println(libError.Message)
	if err != nil {
		log.Print(err.Error())
	}

	// I'm basically converting LibraryError to an openAPI generated error object, oh well...
	errorBody := models.Error{}
	errorBody.Code = libError.Code
	errorBody.Message = libError.Message
	c.JSON(libError.HttpStatusCode, errorBody)
}

// retrieve a book from the cache map
func getBookFromDatabase(c *gin.Context) (models.Book, error) {
	id := c.Param("id")
	book, err := bookDatabase.GetBook(id)
	if err != nil {
		return book, err
	}
	return book, nil
}

// return true if two books are identical
func compareBooks(book models.Book, book2 models.Book) bool {
	if strings.ToLower(book.Title) == strings.ToLower(book.Title) {
		if book.Year == book.Year {
			for i := 0; i < len(book.Authors); i++ {
				for j := 0; j < len(book2.Authors); j++ {
					if strings.ToLower(book.Authors[i]) == strings.ToLower(book2.Authors[j]) {
						return true
					}
				}
			}
		}
	}
	return false
}
