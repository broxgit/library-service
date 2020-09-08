package db

import (
	"github.com/broxgit/library-service/internal/pkg/models"
	"github.com/gocql/gocql"
	"log"
	"os"
	"time"
)

// CassandraDatabase implements the Database interface
type CassandraDatabase struct {
	session *gocql.Session
}

// create a new CassandraDatabase 'object'
func newCassandraDb(keySpace string) (*CassandraDatabase, error) {
	cassandraHost, found := os.LookupEnv("CASSANDRA_HOSTNAME")
	cassandraUser, userFound := os.LookupEnv("CASSANDRA_USERNAME")
	cassandraPass, passFound := os.LookupEnv("CASSANDRA_PASSWORD")
	if !found {
		log.Fatal("Could not find CASSANDRA_HOSTNAME environment value...")
	}
	if !userFound {
		log.Fatal("Could not find CASSANDRA_USERNAME environment value...")
	}
	if !passFound {
		log.Fatal("Could not find CASSANDRA_PASSWORD environment value...")
	}
	cluster := gocql.NewCluster(cassandraHost)
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUser, Password: cassandraPass}
	cluster.Keyspace = keySpace
	cluster.Consistency = gocql.Quorum

	cluster.ReconnectionPolicy = &gocql.ConstantReconnectionPolicy{
		MaxRetries: 10,
		Interval:   30 * time.Second,
	}

	session, err := cluster.CreateSession()
	db := &CassandraDatabase{session: session}

	return db, err
}

// insert new book into database
func (db *CassandraDatabase) CreateBook(book models.Book) error {
	return insertBook(book, db)
}

// get book from database
func (db *CassandraDatabase) GetBook(bookId string) (models.Book, error) {
	var book models.Book

	err := db.session.
		Query(`SELECT id, title, authors, year, version, creationTime, lastUpdateTime, comment FROM books WHERE id = ? LIMIT 1`, bookId).
		Consistency(gocql.One).
		Scan(&book.Id, &book.Title, &book.Authors, &book.Year, &book.Version, &book.CreationTime, &book.LastUpdateTime, &book.Comment)

	return book, err
}

// update book in database
func (db *CassandraDatabase) UpdateBook(book models.Book) error {
	err := insertBook(book, db)
	if err != nil {
		return err
	}
	return nil
}

// delete book from database
func (db *CassandraDatabase) DeleteBook(bookId string) error {
	err := db.session.
		Query(`DELETE FROM books WHERE id = ?`, bookId).
		Exec()
	return err
}

// get all books from database
func (db *CassandraDatabase) ListBooks() ([]models.Book, error) {
	var books []models.Book
	var book models.Book
	iter := db.session.Query(`SELECT id, title, authors, year, version, creationTime, lastUpdateTime, comment FROM books`).Iter()
	for iter.Scan(&book.Id, &book.Title, &book.Authors, &book.Year, &book.Version, &book.CreationTime, &book.LastUpdateTime, &book.Comment) {
		books = append(books, book)
	}
	return books, nil
}

func insertBook(book models.Book, db *CassandraDatabase) error {
	err := db.session.Query(`INSERT INTO books (id, title, authors, year, version, creationTime, lastUpdateTime, comment) values(?, ?, ?, ?, ?, ?, ?, ?)`,
		book.Id, book.Title, book.Authors, book.Year, book.Version, book.CreationTime, book.LastUpdateTime, book.Comment).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
