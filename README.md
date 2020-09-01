# Library API Service

## Project Requirements
Using a language of your choice, write a minimal library API that can perform the following functions:
- List all books in the library
- CRUD operations on a single book

## Technologies Used
GoLang     

Docker

Python 3.6

openAPI 3.0

## Running Locally
The following instructions/commands should be executed in the root of the library-service directory.

Run the following command to run the server locally:
```bash
go run cmd\libraryservice\main.go
```

### Running in Docker 
The following instructions/commands should be executed in the root of the library-service directory.

#### Build the Go Binary:
**Windows Users Only**    
Create an environment variable for Go which will allow Windows users to build binaries for the Linux architecture: `GOOS=linux`    

1.  `go build -o library-service`

#### Run the Docker Commands:
1. `docker build -t library-service .`
2. `docker run -p 8081:8081 library-service`
 
## View the API Definition
On a browser, navigate to http://localhost:8081/swaggerui/

## API Usage Examples with Curl

### Creating a Book
Create a JSON file with payload data (e.g. `payload.json`):
```json
{
    "title": "The Great Gatsby",
    "authors": [
        "F. Scott Fitzgerald"
    ],
    "year": 1925,
    "comment": "The story of the mysteriously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."
}
```

```bash
curl -X POST "http://localhost:8081/library-service/v1/books" -H  "accept: application/json" -H  "Content-Type: application/json" --data-binary @protectData.json
```

### Updating a Book
**Requires:** Id and Version

Modify the payload JSON file created above, or modify the JSON data returned from the API via a GET /books call. 

```bash
curl -X PUT "http://localhost:8081/library-service/v1/books/207e94b8-bc96-446a-b5f0-11c860dae234" -H  "accept: application/json" -H  "If-Match: "fbd34119-3538-4e72-bdcc-3c95b59e8e5b"" -H  "Content-Type: application/json --data-binary @protectData.json"
```

### Deleting a Book
**Requires:** Id and Version

```bash
curl -X DELETE "http://localhost:8081/library-service/v1/books/207e94b8-bc96-446a-b5f0-11c860dae234" -H  "accept: */*" -H  "If-Match: "207e94b8-bc96-446a-b5f0-11c860dae234""
```

### Getting a Single Book
```bash
curl -X GET "http://localhost:8081/library-service/v1/books/207e94b8-bc96-446a-b5f0-11c860dae234" -H  "accept: application/json"
```

### Getting a List of Books
```bash
curl -X GET "http://localhost:8081/library-service/v1/books" -H  "accept: application/json"
```