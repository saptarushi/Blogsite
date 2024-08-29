# Blogsite

[![Go Report Card](https://goreportcard.com/badge/github.com/saptarushi/Blogsite)](https://goreportcard.com/report/github.com/saptarushi/Blogsite)

Blogsite is a blog management REST API built with Golang. The application allows users to register, log in, create, update, delete, and retrieve blog posts. It provides secure user authentication using JWT tokens.

## Features

- User registration and login
- JWT-based authentication and authorization
- CRUD operations for blog posts
- User-specific blog management
- Integration and unit tests
- Dockerized deployment for easy setup and scalability

## Technologies Used

- **Golang**: The primary programming language used for the backend.
- **Gorilla Mux**: A powerful URL router and dispatcher for Golang.
- **GORM**: An ORM library for Golang that simplifies database interactions.
- **JWT**: JSON Web Tokens for secure authentication and authorization.
- **SQLite**: The default database, easily replaceable with other databases.
- **Docker**: Containerization for consistent and portable deployment.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (1.18+)
- [Git](https://git-scm.com/downloads)
- [Docker](https://docs.docker.com/engine/install/) (for containerized deployment)

### Installation

1. **Clone the repository:**

```bash
git clone https://github.com/saptarushi/Blogsite.git
cd Blogsite
```

2.  **Install the dependencies:**

```bash
go mod tidy
```

3.  **Set up the database:**

    By default, the application uses postgreSQL. If you want to use another database, update the connection settings in `main.go` & `database.go`. If you stick with postgreSQL, you can create the database file like this:

```bash
docker-compose up -d
```

4.  **Run the application:**

```bash
go run main.go
```
## Docker Setup
--------------
To run the application using Docker, follow these steps:

1. Build and run the Docker containers:
   
```bash
docker-compose up --build
```
This will start the application and the PostgreSQL database in separate containers.

2. Access the application:
 
Once the containers are up and running, you can interact with the API on http://localhost:8080.

Running Tests
-------------

The application includes both unit and integration tests to ensure code reliability.

To run all tests:

```bash
go test ./...
```
To run specific tests, navigate to the corresponding directory and run:

```bash
go test (directory)
```
## API Documentation
The API endpoints are documented in a Postman collection. You can view and interact with the API using this Postman collection.

Usage:
-------------

Once the application is running, you can interact with the API using tools like Postman or cURL.

Register a User:

```bash
POST /api/register
```
Login a User:

```bash
POST /api/login
```
Create a Blog:
```bash
POST /api/user/blog
```
Update a Blog:
```bash
PUT /api/user/blog/{id}
```
Get a Blog:
```bash
GET /api/blog/{id}
```
Delete a Blog:
```bash
DELETE /api/user/blog/{id}
```
For detailed API usage, refer to the [Postman collection](https://documenter.getpostman.com/view/36157146/2sAXjJ7tN4).

## Adherence to Go Best Practices
------------------------------

This project follows Go best practices, including:

-   **Modular Code Structure**: The project is organized into clearly defined packages (handlers, middlewares, models, etc.) to maintain a clean and scalable codebase.
-   **Error Handling**: Consistent error handling is implemented throughout the codebase. Errors are logged, and meaningful messages are returned to the client.
-   **Use of Context**: The application leverages `context.Context` for passing request-scoped values (e.g., user IDs) across the application, particularly in the authentication middleware.
-   **Go Modules**: The project uses Go modules (`go.mod` and `go.sum`) for dependency management, ensuring reproducible builds and ease of dependency tracking.
-   **Testing**: Comprehensive unit and integration tests are included to ensure code reliability and catch regressions early.

Security Considerations
-----------------------

Security is a priority in the Blogsite application:

-   **JWT Authentication**: The application uses JWT tokens for secure user authentication and authorization. Tokens are validated on every request to protected routes.
-   **Password Hashing**: User passwords are securely hashed using a strong hashing algorithm before being stored in the database, ensuring they are never stored in plain text.
-   **Environment Variables**: Sensitive information like JWT signing keys should be stored in environment variables (.env files) and not hardcoded in the source code.
-   **Input Validation**: Proper input validation and sanitization are applied to prevent common vulnerabilities like SQL injection.


Contributing
------------

Contributions are welcome! Please feel free to submit issues and enhancement requests. If you want to contribute code, fork the repository, create a new branch, and submit a pull request.


Acknowledgements
----------------

-   **Golang** for providing an efficient and fast programming language.
-   **Gorilla Mux** for the routing package.
-   **GORM** for the ORM library.
-   **Postman** for API testing and documentation.
-   **Docker** for containerization.
