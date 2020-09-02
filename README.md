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
go run main.go
```

### Running in Docker
The following instructions/commands should be executed in the root of the library-service directory.

#### Option 1: Pull from DockerHub
```bash
docker pull broxhub/library-service
docker run -p 8081:8081 broxhub/library-service
```

#### Option 2: Build the Docker image locally
```bash
docker build -t library-service .
docker run -p 8081:8081 library-service
```
 
## View the API Definition
On a browser, navigate to http://localhost:8081/swaggerui/

## Running Smoke Tests
Within this repository, there is a Python script that will run some quick tests on the Library Service server.

1. Install the required packages for the script
    ```bash
    pip install -r test/requirements.txt
    ```

2. Run the script
    ```bash
    python test/smoketest.py
    ```

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