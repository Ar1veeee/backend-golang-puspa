```
klinik-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point aplikasi
│
├── internal/
│   ├── handlers/
│   │   ├── auth_handler.go         # Login, register, logout
│   │   ├── admin_handler.go        # Admin endpoints
│   │   ├── terapis_handler.go      # Terapis endpoints
│   │   ├── orangtua_handler.go     # Orang tua endpoints
│   │   ├── workflow_handler.go     # Workflow endpoints
│   │   ├── payment_handler.go      # Payment endpoints
│   │   ├── document_handler.go     # Document endpoints
│   │   ├── notification_handler.go # Notification endpoints
│   │   └── websocket_handler.go    # Real-time updates
│   │
│   ├── services/
│   │   ├── auth_service.go         # Authentication & Authorization
│   │   ├── user_service.go         # User management
│   │   ├── workflow_service.go     # Core workflow logic
│   │   ├── payment_service.go      # Payment per jam
│   │   ├── document_service.go     # Document management
│   │   ├── notification_service.go # Notification system
│   │   ├── reporting_service.go    # Reporting & analytics
│   │   ├── school_service.go       # School integration
│   │   └── file_service.go         # File upload/download
│   │
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── anak_repository.go
│   │   ├── observasi_repository.go
│   │   ├── asesmen_repository.go
│   │   ├── konferensi_repository.go
│   │   ├── intervensi_repository.go
│   │   ├── evaluasi_repository.go
│   │   ├── dokumen_repository.go
│   │   ├── jadwal_repository.go
│   │   ├── pembayaran_repository.go
│   │   └── sekolah_repository.go
│   │
│   └── routes/
│       ├── auth_routes.go         # Authentication routes
│       ├── admin_routes.go        # Admin-specific routes
│       ├── terapis_routes.go      # Terapis-specific routes
│       ├── orangtua_routes.go     # Orang tua routes
│       ├── workflow_routes.go     # Workflow-specific routes
│       ├── payment_routes.go      # Payment routes
│       ├── document_routes.go     # Document management routes
│       └── websocket_routes.go    # WebSocket routes
│
├── shared/
│   ├── models/
│   │   ├── user.go                # Model User (Admin, Terapis, Orangtua)
│   │   ├── anak.go                # Model Anak/Pasien
│   │   ├── observasi.go           # Model Observasi
│   │   ├── asesmen.go             # Model Asesmen
│   │   ├── konferensi.go          # Model Konferensi Kasus
│   │   ├── intervensi.go          # Model Tindakan Intervensi
│   │   ├── evaluasi.go            # Model Evaluasi
│   │   ├── dokumen.go             # Model Dokumen Penyerta
│   │   ├── jadwal.go              # Model Jadwal
│   │   ├── pembayaran.go          # Model Pembayaran
│   │   └── sekolah.go             # Model Integrasi Sekolah
│   │
│   ├── middleware/
│   │   ├── auth_middleware.go      # JWT validation
│   │   ├── role_middleware.go      # Role-based access
│   │   ├── cors_middleware.go      # CORS handling
│   │   ├── rate_limit_middleware.go # Rate limiting
│   │   ├── logging_middleware.go   # Request logging
│   │   └── validation_middleware.go # Input validation
│   │
│   ├── utils/
│   │   ├── jwt.go                  # JWT utilities
│   │   ├── hash.go                 # Password hashing
│   │   ├── validator.go            # Input validation
│   │   ├── file_util.go           # File handling utilities
│   │   ├── email.go               # Email utilities
│   │   ├── sms.go                 # SMS utilities
│   │   └── response.go            # Standard API responses
│   │
│   └── config/
│       ├── config.go              # Configuration management
│       └── database.go            # Database configuration
│
├── pkg/
│   ├── database/
│   │   ├── postgres.go            # PostgreSQL connection driver
│   │   ├── redis.go               # Redis connection driver
│   │   └── migrations/            # Database migration tools
│   │
│   ├── external/
│   │   ├── midtrans/
│   │   │   ├── client.go          # Midtrans API client
│   │   │   └── types.go           # Midtrans data types
│   │   ├── email/
│   │   │   ├── provider.go        # Email service client
│   │   │   └── templates.go       # Email templates
│   │   ├── sms/
│   │   │   └── provider.go        # SMS service client
│   │   ├── storage/
│   │   │   └── cloud_storage.go   # Cloud storage client
│   │   └── school/
│   │       └── integration.go     # School system integration
│   │
│   └── logger/
│       ├── logger.go              # Centralized logging
│       └── formatter.go           # Log formatting
│
├── docs/
│   ├── api/
│   │   └── swagger.yaml           # API documentation
│   └── README.md
│
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── scripts/
│   ├── migrate.go                 # Database migration script
│   ├── seed.go                    # Database seeder
│   └── deploy.sh                  # Deployment script
│
├── docker/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── docker-compose.prod.yml
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```