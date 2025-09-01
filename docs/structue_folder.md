backend-golang/
├── cmd/
│   └── main.go
├── docs/
│   ├── api.md
│   └── deployment.md
├── internal/
│   ├── auth/
│   │   ├── dto/
│   │   │   ├── login_request.go
│   │   │   ├── register_request.go
│   │   │   └── auth_response.go
│   │   ├── handler/
│   │   │   ├── auth_handler.go           # Login, Register, Logout, RefreshToken
│   │   │   └── auth_handler_test.go
│   │   ├── repository/
│   │   │   ├── interfaces.go            # AuthRepository interface
│   │   │   ├── auth_repository.go       # Token management, user lookup
│   │   │   └── auth_repository_test.go
│   │   ├── routes/
│   │   │   └── auth_routes.go
│   │   ├── service/
│   │   │   ├── interfaces.go            # AuthService interface  
│   │   │   ├── auth_service.go          # JWT logic, password validation
│   │   │   └── auth_service_test.go
│   │  
│   │
│   ├── user/
│   │   ├── dto/
│   │   │   ├── user_request.go          # UserCreateRequest, UserUpdateRequest
│   │   │   └── user_response.go         # UserResponse
│   │   ├── handler/
│   │   │   ├── user_handler.go          # CRUD user profile, GetUsers (admin)
│   │   │   └── user_handler_test.go
│   │   ├── repository/
│   │   │   ├── interfaces.go            # UserRepository interface
│   │   │   ├── user_repository.go       # User data operations (all methods)
│   │   │   └── user_repository_test.go
│   │   ├── routes/
│   │   │   └── user_routes.go
│   │   ├── service/
│   │   │   ├── interfaces.go            # UserService interface
│   │   │   ├── user_service.go          # Profile management, user business logic
│   │   │   └── user_service_test.go
│   │   
│   │
│   └── admin/                           # Future expansion
│       ├── dto/
│       ├── handler/
│       ├── repository/
│       ├── routes/
│       ├── service/
│       └── mocks/
├── pkg/
│   └── logger/
│       └── logger.go
├── shared/
│   ├── models/                          # ✅ All database models here
│   │   ├── user.go                     # User model with BeforeCreate hook
│   │   ├── refresh_token.go            # RefreshToken model with User relation
│   │   └── base.go                     # Optional: common fields (ID, CreatedAt, etc.)
│   ├── config/
│   │   ├── config.go
│   │   └── database.go
│   ├── constants/
│   │   ├── errors.go
│   │   ├── status.go
│   │   └── roles.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migrations/
│   │       ├── 001_create_users_table.sql
│   │       ├── 002_create_refresh_tokens_table.sql
│   │       └── 003_create_user_profiles_table.sql
│   ├── helpers/
│   │   ├── hash.go                     # HashPassword function
│   │   ├── jwt.go                      # JWT utilities
│   │   ├── refresh_token.go            # Refresh token utilities  
│   │   └── validator.go                # TranslateErrorMessage function
│   ├── middlewares/
│   │   ├── auth_middleware.go          # JWT authentication
│   │   ├── cors_middleware.go
│   │   ├── logger_middleware.go
│   │   └── rate_limit_middleware.go
│   ├── routes/
│   │   └── routes.go                   # Root route configuration
│   ├── types/
│   │   ├── error_response.go           # ErrorResponse struct
│   │   ├── success_response.go         # SuccessResponse struct
│   │   └── pagination.go
│   └── utils/
│       ├── encryption.go
│       ├── time.go
│       ├── string.go
│       └── file.go
├── tests/
│   ├── integration/
│   │   ├── auth_test.go
│   │   └── user_test.go
│   └── fixtures/
│       ├── users.json
│       └── auth_tokens.json
├── scripts/
│   ├── migrate.sh
│   ├── seed.sh
│   └── build.sh
├── deployments/
│   ├── docker/
│   │   └── Dockerfile
│   └── k8s/
│       ├── deployment.yaml
│       └── service.yaml
├── tmp/
│   ├── build-errors.log
│   ├── main.exe
│   ├── .air.toml
│   └── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md