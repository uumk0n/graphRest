
---

# graphRest

A RESTful API implementation for interacting with graph databases using Go.

## Overview

`graphRest` is a Go-based service that exposes RESTful HTTP endpoints to interact with graph databases. The API supports basic CRUD operations on nodes and relationships, offering the ability to retrieve, create, update, and delete graph elements through HTTP requests.

## Features

- **GET** all nodes with attributes `id` and `label`.
- **GET** a node and all its relationships with connected nodes.
- **POST** to add a new node, relationships, or a graph segment.
- **DELETE** to remove a node, relationships, or a graph segment.
- Secure endpoints using an authorization token.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/uumk0n/graphRest.git
   ```
2. Navigate into the project directory:
   ```bash
   cd graphRest
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Run the application:
```bash
go run main.go
```

By default, the service will start on port `8080`. You can change the port via configuration.

## API Endpoints

- **GET /nodes**: List all nodes.
- **GET /node/{id}**: Retrieve a node and its relationships.
- **POST /node**: Create a new node.
- **DELETE /node/{id}**: Delete a node.

## Configuration

Edit the `config` directory to modify database connections and API settings.

## Tests

Run the tests with:
```bash
go test ./tests
```
