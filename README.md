# Library API Service

- [Library API Service](#library-api-service)
  * [Project Requirements](#project-requirements)
  * [Technologies Used](#technologies-used)
- [Running the Server](#running-the-server)
  * [Running Locally](#running-locally)
  * [Running in Docker](#running-in-docker)
    - [Option 1: Pull from DockerHub](#option-1-pull-from-dockerhub)
    - [Option 2: Build the Docker image locally](#option-2-build-the-docker-image-locally)
  * [Installing with Helm and K8s](#installing-with-helm-and-kubernetes)
    + [Requirements](#requirements)
        + [0. Optional: Modify the values-override.yaml file](#0-optional-modify-the-values-overrideyaml-file)
        + [1. Add the Bitnami Helm Repo](#1-add-the-bitnami-helm-repo)
        + [2. Install the Chart](#2-install-the-chart)
          - [Option 1: Using default values](#option-1-using-default-values)
          - [Option 2: Overriding with values-overrides.yaml](#option-2-overriding-with-values-overridesyaml)
        + [3. Get the service port for the library server](#3-get-the-service-port-for-the-library-server)
  * [View the API Definition](#view-the-api-definition)
  * [Running Smoke Tests](#running-smoke-tests)
      + [Smoke Test Script](#smoke-test-script)
        - [Usage](#usage)
        - [1. Install Requirements](#1-install-requirements)
        - [2. Run the Script](#2-run-the-script)
  * [API Usage Examples with Curl](#api-usage-examples-with-curl)
    + [Creating a Book](#creating-a-book)
    + [Updating a Book](#updating-a-book)
    + [Deleting a Book](#deleting-a-book)
    + [Getting a Single Book](#getting-a-single-book)
    + [Getting a List of Books](#getting-a-list-of-books)
  

## Project Requirements
Using a language of your choice, write a minimal library API that can perform the following functions:
- List all books in the library
- CRUD operations on a single book

## Technologies Used
GoLang     

Docker

Python 3.6

openAPI 3.0

Helm 3

Kubernetes

Cassandra

Bash Scripting

# Running the Server
To run the server without cloning the repo, pull and run the Docker image by following [these instructions.](#option-1-pull-from-dockerhub)

Otherwise, cloning the repo will allow for additional ways to run the server and provides users with a Python script for testing the server functionality.

## Running Locally
The following instructions/commands should be executed in the root of the library-service directory.

Run the following command to run the server locally:
```bash
go run main.go
```

## Running in Docker
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

## Installing with Helm and Kubernetes
A Helm chart for the library service has been included in this repo. The Helm Chart uses Cassandra to manage the storage of library books. 

If using Cassandra (default behavior), the Helm chart will take a few minutes to install while Cassandra starts up and is initialized.

### Requirements
- Helm 3.x
- Kubernetes
- Docker Runtime

### 0. Optional: Modify the values-override.yaml file
The **_values-override.yaml_** file can be modified to fit user specifications (e.g. enabling an ingres service, increasing resources, disabling Cassandra, etc...)

Modifications are not necessary and the chart will work fine with the default values provided.

### 1. Add the Bitnami Helm Repo
```bash
cd build/helm

helm repo add bitnami https://charts.bitnami.com/bitnami
```

### 2. Install the Chart
From the root of the project directory, run the following command:
#### Option 1: Using default values
```bash
helm install library build/helm/library
```

#### Option 2: Overriding with values-overrides.yaml
```bash
helm install library build/helm/library -f values-overrides.yaml
```

### 3. Get the service port for the library server
```bash
kubectl get -o jsonpath="{.spec.ports[0].nodePort}" services library
```

The port returned from the above command will be used to make API calls instead of 8081.
 
## View the API Definition
On a browser, navigate to http://localhost:8081/swaggerui/

## Running Smoke Tests

### Smoke Test Script
Within this repository, there is a Python script that will run some quick tests on the Library Service server.

#### Usage
```bash
usage: smoketest.py [-h] [-p PORT]

Library Service Smoke Tests

optional arguments:
  -h, --help            show this help message and exit
  -p PORT, --port PORT  The port where the library service is hosted
```

#### 1. Install Requirements 
```bash
pip install -r test/requirements.txt
```

#### 2. Run the Script
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