# Backend Golang API Documentation

## Overview

This is a RESTful API built with Go (Golang) using the Gin framework, GORM for database operations, and JWT for authentication. The API provides user management and authentication services with a clean, layered architecture following clean architecture principles.

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
    "field": "error message"
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
  "password": "securepassword123"
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
    "role": "therapist",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
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
    "role": "therapist",
    "accessToken": "jwt_access_token",
    "tokenType": "Bearer",
    "refreshToken": "refresh_token_string",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
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
  "refreshToken": "your_refresh_token"
}
```
- **Response:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "accessToken": "new_jwt_access_token",
    "refreshToken": "new_refresh_token",
    "tokenType": "Bearer",
    "expiresIn": "2024-01-01T00:00:00Z"
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
      "role": "therapist",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
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
  "role": "therapist"
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
    "role": "therapist",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
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
    "role": "therapist",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
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
  "role": "admin"
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
    "role": "admin",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
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

## User Roles

The application supports three user roles defined in `shared/constants/role.go`:

```go
type Role string

const (
    RoleAdmin        Role = "admin"
    RoleUser         Role = "therapist"
    RolePsychiatrist Role = "psychiatrist"
)
```

- **`admin`**: Full access to all endpoints, can create, read, update, and delete users
- **`user`**: Standard user with limited access, can only read their own profile
- **`psychiatrist`**: Specialized role for mental health professionals

Role-based access control is implemented at the middleware level to ensure proper authorization. Users can only access endpoints appropriate to their role level.

## Data Models

### User Model
```go
type User struct {
    Id        string    `json:"id" gorm:"primary_key;type:char(36);"`
    Name      string    `json:"name"`
    Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
    Email     string    `json:"email" gorm:"type:varchar(50);uniqueIndex;not null"`
    Password  string    `json:"password"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt;autoCreateTime"`
    UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt;autoUpdateTime"`
}
```

### Refresh Token Model
```go
type RefreshToken struct {
    Id        string    `json:"id" gorm:"primary_key;type:char(36);"`
    UserId    string    `json:"user_id" gorm:"type:char(36);not null;index"`
    Token     string    `json:"token" gorm:"not null"`
    ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
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
- `name`: Required, string, minimum 3 characters, maximum 100 characters
- `username`: Required, unique, string, minimum 3 characters, maximum 50 characters, alphanumeric only
- `email`: Required, unique, valid email format
- `password`: Required, string, minimum 8 characters
- `role`: Required for user creation, optional for update, maximum 15 characters (admin, user, psychiatrist)

### Authentication
- `username`: Required, string, minimum 3 characters, maximum 50 characters, alphanumeric only
- `password`: Required, string
- `refreshToken`: Required, string

## Middleware

### Authentication Middleware
- JWT token validation
- User context injection
- Token expiration checking

### Logger Middleware
- Request/response logging
- Performance metrics
- Error tracking

### Rate Limiting Middleware
- Request rate limiting per IP
- Configurable limits and windows
- Graceful handling of exceeded limits

## CORS Configuration

The API supports CORS with the following configuration:
- **Allowed Origins:** Configurable via environment
- **Allowed Methods:** GET, POST, PUT, PATCH, DELETE, OPTIONS
- **Allowed Headers:** Origin, Content-Type, Authorization, Accept
- **Exposed Headers:** Content-Length, X-Total-Count
- **Credentials:** true

## Environment Variables

The application uses the following environment variables (configured via `.env` file):

### Application
- `APP_PORT`: Server port (default: 3000)
- `GIN_MODE`: Gin mode (debug/release, default: debug)

### Database
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 3306)
- `DB_USER`: Database username
- `DB_PASS`: Database password
- `DB_NAME`: Database name
- `DB_ROOT_PASSWORD`: Database root password

### JWT
- `JWT_SECRET`: JWT signing secret
- `JWT_EXPIRY`: JWT expiration time (default: 24h)
- `REFRESH_TOKEN_EXPIRY`: Refresh token expiration time (default: 168h)

## Dependencies

### Core Dependencies
- **Gin:** Web framework for HTTP routing
- **GORM:** ORM for database operations
- **JWT:** Authentication tokens (github.com/golang-jwt/jwt/v5)
- **MySQL:** Database driver (gorm.io/driver/mysql)
- **UUID:** Unique identifier generation (github.com/google/uuid)
- **Bcrypt:** Password hashing (golang.org/x/crypto/bcrypt)

### Development Dependencies
- **Godotenv:** Environment variable loading
- **CORS:** Cross-origin resource sharing
- **Air:** Hot reload for development
- **Golang-migrate:** Database migrations

## Getting Started

### Prerequisites
- Go 1.21 or higher
- MySQL 8.0 or higher
- Docker and Docker Compose (optional)

### Local Development
1. **Clone the repository**
2. **Install dependencies:** `go mod download`
3. **Set up environment variables** in `.env` file
4. **Run database migrations:** `go run shared/database/migrations.go`
5. **Run the application:** `go run cmd/main.go`

### Docker Development
1. **Set up environment variables** in `.env` file
2. **Build and run:** `docker-compose up -d`
3. **View logs:** `docker-compose logs -f app`

### Database Setup
1. **Create database:** MySQL database with the specified name
2. **Run migrations:** Execute migration files in `shared/database/migrations/`
3. **Verify connection:** Check database connectivity

## Architecture

The application follows a clean, layered architecture:

```
├── cmd/                    # Application entry point
│   └── main.go            # Main application file
├── internal/               # Business logic (private packages)
│   ├── auth/               # Authentication module
│   │   ├── dto/            # Data Transfer Objects
│   │   ├── errors/         # Domain-specific errors
│   │   ├── handler/        # HTTP request handlers
│   │   ├── repository/     # Data access layer
│   │   ├── routes/         # Route definitions
│   │   └── service/        # Business logic layer
│   └── user/               # User management module
│       ├── dto/            # Data Transfer Objects
│       ├── errors/         # Domain-specific errors
│       ├── handler/        # HTTP request handlers
│       ├── repository/     # Data access layer
│       ├── routes/         # Route definitions
│       └── service/        # Business logic layer
├── shared/                 # Shared utilities and configurations
│   ├── config/             # Configuration management
│   ├── constants/          # Application constants
│   ├── database/           # Database connection and migrations
│   ├── errors/             # Common error types
│   ├── helpers/            # Utility functions
│   ├── logger/             # Logging configuration
│   ├── middlewares/        # HTTP middlewares
│   ├── models/             # Data models
│   ├── routes/             # Route setup and middleware
│   └── types/              # Common types and interfaces
├── docs/                   # Documentation
├── docker-compose.yaml     # Docker services configuration
├── Dockerfile              # Application container definition
└── go.mod                  # Go module dependencies
```

### Module Pattern
Each module follows the same pattern:
- **DTOs:** Data Transfer Objects for request/response validation
- **Handlers:** HTTP request handlers (controller layer)
- **Services:** Business logic layer (use case layer)
- **Repositories:** Data access layer (data layer)
- **Routes:** Route definitions and middleware setup
- **Errors:** Domain-specific error types

### Dependency Flow
```
HTTP Request → Handler → Service → Repository → Database
     ↓
HTTP Response ← Handler ← Service ← Repository ← Database
```

## Security Features

- **JWT-based authentication** with configurable expiration
- **Password hashing** using bcrypt with salt
- **CORS protection** with configurable origins
- **Input validation** using custom validators
- **Rate limiting** to prevent abuse
- **Graceful shutdown** handling
- **Database connection** pooling and management
- **Secure headers** and middleware
- **Token refresh** mechanism for better security

## Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific module tests
go test ./internal/auth/...
go test ./internal/therapist/...
```

### Test Structure
- Unit tests for services and repositories
- Integration tests for handlers
- Mock implementations for external dependencies

## Deployment

### Production Considerations
- Set `GIN_MODE=release` for production
- Use strong JWT secrets
- Configure proper CORS origins
- Set up database connection pooling
- Enable logging and monitoring
- Use HTTPS in production
- Set up proper firewall rules

### Environment-Specific Configs
- Development: Local database, debug logging
- Staging: Staging database, info logging
- Production: Production database, error logging only

## Monitoring and Logging

### Logging Levels
- **DEBUG:** Detailed information for debugging
- **INFO:** General information about application flow
- **WARN:** Warning messages for potential issues
- **ERROR:** Error messages for failed operations

### Metrics
- Request/response times
- Database query performance
- Error rates and types
- Authentication success/failure rates

## Troubleshooting

### Common Issues
1. **Database Connection Failed:** Check database credentials and network
2. **JWT Token Invalid:** Verify JWT secret and token expiration
3. **CORS Errors:** Check allowed origins configuration
4. **Rate Limiting:** Adjust rate limit configuration if needed

### Debug Mode
Enable debug mode by setting `GIN_MODE=debug` to get detailed error information and request logging.

## Contributing

1. Follow the existing code structure and patterns
2. Add tests for new functionality
3. Update documentation for API changes
4. Follow Go coding standards and best practices
5. Use meaningful commit messages

## License

This project is licensed under the MIT License.

