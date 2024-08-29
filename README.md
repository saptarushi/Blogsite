# Blogsite

[![Go Report Card](https://goreportcard.com/badge/github.com/saptarushi/Blogsite)](https://goreportcard.com/report/github.com/saptarushi/Blogsite)

Blogsite is a blog management REST API built with Golang. The application allows users to register, log in, create, update, delete, and retrieve blog posts. It provides secure user authentication using JWT tokens.

## Features

- User registration and login
- JWT-based authentication and authorization
- CRUD operations for blog posts
- User-specific blog management
- Integration and unit tests

## Technologies Used

- **Golang**: The primary programming language used for the backend.
- **Gorilla Mux**: A powerful URL router and dispatcher for Golang.
- **GORM**: An ORM library for Golang that simplifies database interactions.
- **JWT**: JSON Web Tokens for secure authentication and authorization.
- **SQLite**: The default database, easily replaceable with other databases.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (1.18+)
- [Git](https://git-scm.com/downloads)
- A database (SQLite is used by default, but you can configure it to use MySQL, PostgreSQL, etc.)

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/saptarushi/Blogsite.git
   cd Blogsite
