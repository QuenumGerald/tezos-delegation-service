# Delegation Application in Go

## Table of Contents

1. [Introduction](#introduction)
2. [Dependencies](#dependencies)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Usage](#usage)
6. [API](#api)
7. [Code Structure](#code-structure)
8. [Error Handling](#error-handling)
9. [Future Improvements](#future-improvements)

## Introduction

This application fetches and stores delegation data from Tezos blockchain. It uses Go as the programming language, SQLite3 for database storage, and Gorilla Mux for HTTP routing.

## Dependencies

- **Go**: The programming language used.
- **SQLite3**: The database for storing delegation data.
- **Gorilla Mux**: A HTTP router for Go.

To install Gorilla Mux:
go get -u github.com/gorilla/mux


## Installation

1. Install Go and SQLite3.
2. Clone the repository to your local machine.
3. Navigate to the directory where `main.go` is located.
4. Run `go run main.go`.

## Configuration

The SQLite database (`delegations.db`) is located in the same directory as the `main.go` file. If you need to change the directory, modify the database path in the `initDB()` function. The database is configured to have a table named `delegations` with the following columns:

- `timestamp`: Stores the time of the delegation.
- `amount`: Stores the amount involved in the delegation.
- `delegator`: Stores the address of the delegator.
- `block`: Stores the block hash.

## Usage

1. Run the application: `go run main.go`
2. Access `http://localhost:8000/xtz/delegations` to see a list of delegations.

## API

- **Endpoint**: `GET /xtz/delegations`
    - **Description**: Fetches the delegations from the database and returns them as a JSON array.

### Code Structure

#### Types

- `Sender`: Struct to hold the sender's address.
  ```go
  type Sender struct {
      Address string `json:"address"`
  }

- `Delegation`: Struct to hold the delegation data.
type Delegation struct {
    Timestamp string `json:"timestamp"`
    Amount    int64  `json:"amount"`
    Delegator string `json:"delegator"`
    Block     string `json:"block"`
}

### Functions

- `initDB()`: Initializes the SQLite database and returns the `*sql.DB` instance. It creates the `delegations` table if it doesn't exist.
  
- `fetchDelegations(db *sql.DB)`: Periodically fetches data from the Tezos API and inserts new delegations into the database.

- `getDelegations(w http.ResponseWriter, r *http.Request, db *sql.DB)`: Retrieves delegations from the SQLite database and sends them as a JSON-encoded HTTP response.

## Error Handling

- API fetch errors and database errors are logged but do not stop the application.
- JSON marshaling and unmarshaling errors are also logged.

## Future Improvements

- Add pagination to the `GET /xtz/delegations` API.
- Add more robust error handling and logging.
