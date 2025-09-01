# Backend Golang API Documentation

## Overview

This is a RESTful API built with Go (Golang) using the Gin framework, GORM for database operations, and JWT for authentication. The API provides user management and authentication services with a clean, layered architecture.

## Base URL

```
http://localhost:3000/api/v1
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {}
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "errors": {
    "field": "errors message"
  }
}
```

## API Endpoints

### Authentication Endpoints

#### 1. User Registration
- **URL:** `POST /auth/register`
- **Description:** Register a new user account
- **Authentication:** Not required
- **Request Body:**
```json
{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "business_name": "John's Business"
}
```
- **Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid-string",
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "business_name": "John's Business",
    "create_at": "2024-01-01T00:00:00Z",
    "update_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 2. User Login
- **URL:** `POST /auth/login`
- **Description:** Authenticate user and receive access token
- **Authentication:** Not required
- **Request Body:**
```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```
- **Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "id": "uuid-string",
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "business_name": "John's Business",
    "create_at": "2024-01-01T00:00:00Z",
    "update_at": "2024-01-01T00:00:00Z",
    "token": "jwt_access_token"
  }
}
```

#### 3. Refresh Token
- **URL:** `POST /auth/refresh`
- **Description:** Get a new access token using refresh token
- **Authentication:** Not required
- **Request Body:**
```json
{
  "refresh_token": "your_refresh_token"
}
```
- **Response:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "new_jwt_access_token"
  }
}
```

#### 4. User Logout
- **URL:** `POST /auth/logout`
- **Description:** Logout user and invalidate tokens
- **Authentication:** Required
- **Request Body:** Empty
- **Response:**
```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

### User Management Endpoints

All user management endpoints require authentication.

#### 1. Get All Users
- **URL:** `GET /users`
- **Description:** Retrieve a list of all users
- **Authentication:** Required
- **Query Parameters:** None
- **Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "uuid-string",
      "name": "John Doe",
      "username": "johndoe",
      "email": "john@example.com",
      "business_name": "John's Business",
      "create_at": "2024-01-01T00:00:00Z",
      "update_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 2. Create User
- **URL:** `POST /users`
- **Description:** Create a new user (admin only)
- **Authentication:** Required
- **Request Body:**
```json
{
  "name": "Jane Doe",
  "username": "janedoe",
  "email": "jane@example.com",
  "password": "securepassword123",
  "business_name": "Jane's Business"
}
```
- **Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "uuid-string",
    "name": "Jane Doe",
    "username": "janedoe",
    "email": "jane@example.com",
    "business_name": "Jane's Business",
    "create_at": "2024-01-01T00:00:00Z",
    "update_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3. Get User by ID
- **URL:** `GET /users/{id}`
- **Description:** Retrieve a specific user by ID
- **Authentication:** Required
- **Path Parameters:**
  - `id`: User UUID
- **Response:**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "uuid-string",
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "business_name": "John's Business",
    "create_at": "2024-01-01T00:00:00Z",
    "update_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4. Update User
- **URL:** `PUT /users/{id}`
- **Description:** Update an existing user
- **Authentication:** Required
- **Path Parameters:**
  - `id`: User UUID
- **Request Body:**
```json
{
  "name": "John Updated",
  "username": "johnupdated",
  "email": "john.updated@example.com",
  "password": "newpassword123",
  "business_name": "John's Updated Business"
}
```
- **Response:**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": "uuid-string",
    "name": "John Updated",
    "username": "johnupdated",
    "email": "john.updated@example.com",
    "business_name": "John's Updated Business",
    "create_at": "2024-01-01T00:00:00Z",
    "update_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 5. Delete User
- **URL:** `DELETE /users/{id}`
- **Description:** Delete a user
- **Authentication:** Required
- **Path Parameters:**
  - `id`: User UUID
- **Response:**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

## Data Models

### User Model
```go
type User struct {
    Id           string    `json:"id" gorm:"primary_key;type:char(36);"`
    Name         string    `json:"name"`
    Username     string    `json:"username" gorm:"unique;not null"`
    Email        string    `json:"email" gorm:"unique;not null"`
    Password     string    `json:"password"`
    BusinessName string    `json:"business_name"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
```

### Refresh Token Model
```go
type RefreshToken struct {
    Id        string    `json:"id" gorm:"primary_key;type:char(36);"`
    UserId    string    `json:"userId" gorm:"type:char(36);not null;index"`
    Token     string    `json:"token" gorm:"not null"`
    ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
    CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
    Revoked   bool      `json:"revoked" gorm:"default:false"`
    User      User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
```

## Error Codes

The API uses standard HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Validation Error
- `500` - Internal Server Error

## Validation Rules

### User Registration/Update
- `name`: Required, string
- `username`: Required, unique, string
- `email`: Required, unique, valid email format
- `password`: Required, string
- `business_name`: Required, string

### Authentication
- `username`: Required, string
- `password`: Required, string
- `refresh_token`: Required, string

## CORS Configuration

The API supports CORS with the following configuration:
- **Allowed Origins:** `http://localhost:3000`
- **Allowed Methods:** GET, POST, PUT, PATCH, DELETE, OPTIONS
- **Allowed Headers:** Origin, Content-Type, Authorization
- **Exposed Headers:** Content-Length

## Environment Variables

The application uses the following environment variables (configured via `.env` file):

- `APP_PORT`: Server port (default: 3000)
- Database configuration variables (configured in your environment)

## Dependencies

### Core Dependencies
- **Gin:** Web framework
- **GORM:** ORM for database operations
- **JWT:** Authentication tokens
- **MySQL:** Database driver
- **UUID:** Unique identifier generation

### Development Dependencies
- **Godotenv:** Environment variable loading
- **CORS:** Cross-origin resource sharing

## Getting Started

1. **Clone the repository**
2. **Install dependencies:** `go mod download`
3. **Set up environment variables** in `.env` file
4. **Run the application:** `go run cmd/main.go`

## Architecture

The application follows a clean, layered architecture:

```
├── cmd/           # Application entry point
├── internal/      # Business logic
│   ├── auth/      # Authentication module
│   └── user/      # User management module
├── shared/        # Shared utilities and configurations
│   ├── config/    # Configuration management
│   ├── database/  # Database connection
│   ├── helpers/   # Utility functions
│   ├── middlewares/ # HTTP middlewares
│   ├── models/    # Data models
│   ├── routes/    # Route setup
│   └── types/     # Common types
└── docs/          # Documentation
```

Each module follows the pattern:
- **DTOs:** Data Transfer Objects for requests/responses
- **Handlers:** HTTP request handlers
- **Services:** Business logic layer
- **Repositories:** Data access layer
- **Routes:** Route definitions

## Security Features

- JWT-based authentication
- Password hashing
- CORS protection
- Input validation
- Graceful shutdown handling
- Database connection management

