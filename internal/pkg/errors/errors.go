package errors

import "fmt"

type LibraryError struct {
	Code           int32
	Message        string
	HttpStatusCode int
}

func JSON_PARSE_ERROR() LibraryError {
	return LibraryError{1000, "An error was encountered when parsing the JSON request body.", 400}
}

func BOOK_NOT_FOUND(id string) LibraryError {
	return LibraryError{1001, fmt.Sprintf("Book with specified id: %v was not found", id), 400}
}

func BOOK_ALREADY_EXISTS(id string) LibraryError {
	return LibraryError{1002, fmt.Sprintf("Book already exists with id: %v", id), 400}
}
