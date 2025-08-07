# GO-X1 REST API

A simple REST API built with Go and Fiber for learning purposes.

## Features

- **Full CRUD Operations**: Create, Read, Update, Delete users
- **Input Validation**: Request validation using go-playground/validator
- **Database Support**: MySQL integration using GORM
- **Environment Configuration**: Configurable via environment variables
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **CORS Support**: Cross-Origin Resource Sharing enabled
- **Logging**: Request logging middleware
- **UUID Generation**: Utility endpoint for generating UUIDs

## Tech Stack

- **Go** - Programming language
- **Fiber v2** - Web framework
- **GORM** - ORM for database operations
- **MySQL** - Database
- **godotenv** - Environment variable loading
- **validator/v10** - Input validation

## Project Structure

```
.
├── auth/           # Authentication utilities
├── connectDB/      # Database connection
├── models/         # Data models
├── main.go         # Main application file
├── go.mod          # Go module file
├── .env           # Environment variables
└── .env.example   # Environment variables example
```

## Setup and Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/SADBOYISME/GO-X1.git
   cd GO-X1
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env file with your database credentials
   ```

4. **Set up MySQL database**
   ```sql
   CREATE DATABASE go_api_db;
   ```

5. **Build and run**
   ```bash
   go build
   ./GO-X1
   ```

## Environment Variables

| Variable     | Description              | Default          |
|--------------|--------------------------|------------------|
| PORT         | Server port              | 8080             |
| DB_HOST      | Database host            | localhost        |
| DB_PORT      | Database port            | 3306             |
| DB_USER      | Database username        | root             |
| DB_PASSWORD  | Database password        |                  |
| DB_NAME      | Database name            | go_api_db        |
| JWT_SECRET   | Secret key for JWT       | your-secret-key  |

## API Endpoints

### Health Check
- **GET** `/` - API health check
- **GET** `/health` - Detailed health check with database status

### Utilities
- **GET** `/uuid` - Generate a new UUID

### Auth
- **POST** `/api/v1/auth/login` - Login to get a JWT token

### Users (CRUD)
- **POST** `/api/v1/users` - Create a new user
- **GET** `/api/v1/users` - Get all users (protected)
- **GET** `/api/v1/users/:id` - Get user by ID (protected)
- **PUT** `/api/v1/users/:id` - Update user by ID (protected)
- **DELETE** `/api/v1/users/:id` - Delete user by ID (protected)

## API Usage Examples

### Get API Health Status
```bash
curl http://localhost:8080/health
```

### Generate UUID
```bash
curl http://localhost:8080/uuid
```

### Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

### Create a New User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johnupdated",
    "email": "johnupdated@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Response Format

All API responses follow a standard format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "error": null
}
```

### Success Response
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Username must be at least 3 characters long"
}
```

## User Model

### User Fields
- `id` (uint) - Unique identifier (auto-generated)
- `username` (string) - Username (3-50 characters, unique)
- `email` (string) - Email address (valid email format, unique)
- `password` (string) - Password (minimum 6 characters)
- `created_at` (timestamp) - Creation time (auto-generated)
- `updated_at` (timestamp) - Last update time (auto-generated)

### Validation Rules
- **Username**: Required, 3-50 characters, unique
- **Email**: Required, valid email format, unique
- **Password**: Required, minimum 6 characters

## Running with Docker

1.  **Ensure Docker is installed** on your system.

2.  **Create a `.env` file** from the example:
    ```bash
    cp .env.example .env
    ```
    *Note: The default `DB_HOST` is already set to `db` to connect to the Docker container.*

3.  **Build and run the containers:**
    ```bash
    docker-compose up -d --build
    ```
    This command will build the Go application image, pull the MySQL image, and start both containers in detached mode.

4.  **Check the logs:**
    ```bash
    docker-compose logs -f app
    ```

5.  **Stop the containers:**
    ```bash
    docker-compose down
    ```

## Development

### Building the Application
```bash
go build
```

### Running in Development Mode
```bash
go run main.go
```

### Running Tests (when available)
```bash
go test ./...
```

## Database Schema

The application uses GORM for database operations. The database schema is automatically migrated when the application starts.

### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Error Handling

The API includes comprehensive error handling:

- **400 Bad Request** - Invalid input or malformed JSON
- **404 Not Found** - Resource not found
- **422 Unprocessable Entity** - Validation errors
- **500 Internal Server Error** - Server errors
- **503 Service Unavailable** - Database unavailable

## Security

This project implements the following security features:

- **Password Hashing**: Passwords are hashed using bcrypt before being stored in the database.
- **JWT Authentication**: Protected routes require a valid JSON Web Token (JWT).

For production use, consider implementing:

- **Authorization**: Role-based access control (RBAC) to manage user permissions.
- **Rate Limiting**: To prevent API abuse.
- **HTTPS**: To secure communication between the client and the server.
- **Input Sanitization**: For additional protection against injection attacks.

## Contributing

This is a learning project. Feel free to:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

This project is for educational purposes.