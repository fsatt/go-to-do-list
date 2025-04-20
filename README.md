# go-to-do-list
Proof of concept: A simple RESTful API for managing a to-do list, built with Go and PostgreSQL.

## Features

- Add new to-do items
- View all to-do items
- Delete to-do items by ID

## Prerequisites

Before running the application, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.24 or later)
- [PostgreSQL](https://www.postgresql.org/download/)

## Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/fsatt/go-to-do-list.git
   cd go-to-do-list
   ```

2. **Initialize the database**:
   - Create a PostgreSQL database named `todolist`:
     ```bash
     psql -U <your-username> -h localhost -p <your-port> -c "CREATE DATABASE \"todolist\";"
     ```
   - Replace `<your-username>` with your PostgreSQL username, and `<your-port>` with the port your PostgreSQL server is running on.

3. **Update database credentials**:
   - Open main.go and update the `dsn` variable with your PostgreSQL credentials:
     ```go
     dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
         "localhost", "your-port", "your-username", "your-password", "todolist")
     ```

4. **Install dependencies**:
   ```bash
   go mod tidy
   ```

5. **Run the application**:
   ```bash
   go run .
   ```

6. **Access the API**:
   The API will be available at `http://localhost:8080`.

## API Endpoints

### 1. Add a New To-Do Item
**POST** `/todos`

- **Request Body**:
  ```json
  {
    "task": "Your task description"
  }
  ```

- **Response**:
  - `201 Created`:
    ```json
    {
      "message": "Item added"
    }
    ```
  - `400 Bad Request` if the input is invalid.
  - `500 Internal Server Error` if adding item fails.

---

### 2. View All To-Do Items
**GET** `/todos`

- **Response**:
  - `200 OK`:
    ```json
    [
      {
        "id": 1,
        "task": "Your task description"
      }
    ]
    ```
  - `500 Internal Server Error` if fetching items fails.

---

### 3. Delete a To-Do Item by ID
**DELETE** `/todos/:id`

- **Response**:
  - `200 OK`:
    ```json
    {
      "message": "Item successfully deleted"
    }
    ```
  - `404 Not Found` if the item does not exist.
  - `500 Internal Server Error` if deletion fails.

## Project Structure

```
go-to-do-list/
├── main.go       # Main application file
├── go.mod        # Go module file
└── README.md     # Project documentation
```
