package app

import (
	"github.com/broxgit/library-service/internal/pkg/errors"
	"github.com/broxgit/library-service/internal/pkg/models"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

func init() {
	bookCache = make(map[string]*models.Book)
}

var bookCache map[string]*models.Book

// CreateBook - Create a book
func CreateBook(c *gin.Context) {
	book := models.Book{}

	err := c.BindJSON(&book)
	if err != nil {
		handleError(c, err, errors.JSON_PARSE_ERROR())
		return
	}

	// check if book already exists
	for _, _book := range bookCache {
		if compareBooks(&book, _book) {
			handleError(c, nil, errors.BOOK_ALREADY_EXISTS(_book.Id))
			return
		}
	}

	// create book metadata
	book.CreationTime = time.Now()
	book.LastUpdateTime = time.Now()
	book.Id = guuid.New().String()
	book.Version = guuid.New().String()
	bookCache[book.Id] = &book
	c.JSON(http.StatusCreated, book)
}

// DeleteBookById - Deletes a book for a given id.
func DeleteBookById(c *gin.Context) {
	// delete is a no-op if the element doesn't exist,
	// so we shouldn't receive an error

	book := getBookFromCache(c)
	if book != nil {
		// verify that If-Match equals current version
		if book.Version == c.GetHeader("If-Match") {
			delete(bookCache, c.Param("id"))
			c.JSON(http.StatusNoContent, gin.H{})
		} else {
			handleError(c, nil, errors.INVALID_IF_MATCH(book.Version, c.GetHeader("If-Match")))
			return
		}
	} else {
		// book didn't exist in the cache, return an error
		handleError(c, nil, errors.BOOK_NOT_FOUND(c.Param("id")))
		return
	}
}

// GetBookById - Returns a book for a given id.
func GetBookById(c *gin.Context) {
	book := getBookFromCache(c)
	if book != nil {
		c.JSON(http.StatusOK, book)
	} else {
		// book didn't exist in the cache, return an error
		handleError(c, nil, errors.BOOK_NOT_FOUND(c.Param("id")))
		return
	}
}

// ListBooks - List all books
func ListBooks(c *gin.Context) {
	var books []models.Book

	// initialize empty array in the case that we don't have any cached books
	books = make([]models.Book, 0)

	// build the 'books' array
	for _, book := range bookCache {
		books = append(books, *book)
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

	bookCheck := getBookFromCache(c)

	// if book doesn't exist in cache, return an error
	if bookCheck == nil {
		handleError(c, nil, errors.BOOK_NOT_FOUND(c.Param("id")))
		return
	}

	// verify that If-Match equals current version
	ifMatch := c.GetHeader("If-Match")
	if bookCheck.Version == ifMatch {
		// update book metadata
		book.LastUpdateTime = time.Now()
		book.Version = guuid.New().String()
		bookCache[book.Id] = &book

		c.JSON(http.StatusOK, book)
	} else {
		handleError(c, nil, errors.INVALID_IF_MATCH(book.Version, ifMatch))
		return
	}
}

func handleError(c *gin.Context, err error, libError errors.LibraryError) {
	log.Println(libError.Message)
	if err != nil {
		log.Fatal(err.Error())
	}

	// I'm basically converting LibraryError to an openAPI generated error object, oh well...
	errorBody := models.Error{}
	errorBody.Code = libError.Code
	errorBody.Message = libError.Message
	c.JSON(libError.HttpStatusCode, errorBody)
}

// retrieve a book from the cache map
func getBookFromCache(c *gin.Context) *models.Book {
	id := c.Param("id")
	book, found := bookCache[id]
	if !found {
		return nil
	}
	return book
}

// return true if two books are identical
func compareBooks(book *models.Book, book2 *models.Book) bool {
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
