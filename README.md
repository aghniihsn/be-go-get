# Go Backend API - Cinema Ticket Booking System

A comprehensive RESTful API backend service for cinema ticket booking system built with Go, featuring MongoDB ObjectID integration, JWT authentication, and extensive testing coverage.

## 🏗️ Project Structure

```
be-go-get/
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module checksums
├── Makefile                    # Build and deployment scripts
│
├── config/
│   └── database.go            # Database configuration and connection
│
├── controllers/               # HTTP request handlers
│   ├── filmController.go      # Film management endpoints
│   ├── jadwalController.go    # Schedule management endpoints
│   ├── pembayaranController.go # Payment processing endpoints
│   ├── tiketController.go     # Ticket booking endpoints
│   └── userController.go      # User authentication endpoints
│
├── middlewares/
│   └── middleware.go          # JWT authentication middleware
│
├── models/
│   └── struct.go             # Data structures with MongoDB ObjectID
│
├── routes/
│   └── router.go             # API route definitions
│
├── tests/
│   ├── unit/                 # Unit tests for core functionality
│   │   ├── objectid_test.go  # ObjectID validation and generation tests
│   │   ├── password_test.go  # Password hashing security tests
│   │   ├── jwt_test.go       # JWT token authentication tests
│   │   └── controller_logic_test.go # Business logic validation tests
│   │
│   ├── integration/          # Integration tests for API endpoints
│   │   ├── auth_controller_test.go    # Authentication flow tests
│   │   ├── crud_controller_test.go    # CRUD operations tests
│   │   ├── payment_controller_test.go # Payment processing tests
│   │   └── api_test.go               # End-to-end API tests
│   │
│   └── mocks/
│       └── test_data.go      # Test helper functions and mock data
│
├── docs/                     # Project documentation
│   ├── TESTING_PLAN.md       # Comprehensive testing strategy
│   ├── TEST_REPORT.md        # Testing results and analysis
│   ├── CONTROLLER_TEST_REPORT.md # Controller-specific test documentation
│   ├── README_POSTMAN.md     # API usage examples and Postman collection
│   ├── docs.go              # Swagger documentation generator
│   ├── swagger.json         # OpenAPI specification
│   └── swagger.yaml         # OpenAPI specification (YAML)
│
└── scripts/
    └── add-swagger.sh        # Swagger documentation setup script
```

## 🚀 Features

- **MongoDB ObjectID Integration**: Complete migration to MongoDB's ObjectID system for all entities
- **JWT Authentication**: Secure token-based authentication with ObjectID support
- **Comprehensive Testing**: 15+ test cases covering unit and integration testing
- **API Documentation**: Complete Swagger/OpenAPI documentation
- **Clean Architecture**: Well-organized folder structure following Go best practices
- **Security**: Password hashing with bcrypt and secure JWT implementation

## 🛠️ Technologies

- **Language**: Go 1.21+
- **Database**: MongoDB with primitive.ObjectID
- **Authentication**: JWT (JSON Web Tokens)
- **Testing**: Go testing framework with testify assertions
- **Documentation**: Swagger/OpenAPI 3.0
- **Security**: bcrypt password hashing

## 📋 Prerequisites

- Go 1.21 or higher
- MongoDB instance
- Git

## 🔧 Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd be-go-get
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   # Set your MongoDB connection string
   export MONGODB_URI="mongodb://localhost:27017"
   export DATABASE_NAME="cinema_booking"
   export JWT_SECRET="your-secret-key"
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

## 🧪 Testing

The project includes comprehensive testing coverage with both unit and integration tests.

### Run All Tests
```bash
# Run all tests
go test ./tests/...

# Run with verbose output
go test -v ./tests/...
```

### Run Specific Test Suites
```bash
# Unit tests only
go test -v ./tests/unit/

# Integration tests only
go test -v ./tests/integration/
```

### Test Coverage
```bash
# Generate test coverage report
go test -cover ./tests/...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

## 📚 API Documentation

### Access Swagger Documentation
Once the server is running, access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Core Endpoints

#### Authentication
- `POST /api/register` - User registration
- `POST /api/login` - User login

#### Films
- `GET /api/films` - Get all films
- `POST /api/films` - Create new film (admin)
- `GET /api/films/{id}` - Get film by ObjectID
- `PUT /api/films/{id}` - Update film (admin)
- `DELETE /api/films/{id}` - Delete film (admin)

#### Schedules (Jadwal)
- `GET /api/jadwal` - Get all schedules
- `POST /api/jadwal` - Create new schedule (admin)
- `GET /api/jadwal/{id}` - Get schedule by ObjectID

#### Tickets
- `POST /api/tiket` - Book ticket
- `GET /api/tiket/user/{user_id}` - Get user tickets

#### Payments
- `POST /api/pembayaran` - Process payment
- `GET /api/pembayaran/{id}` - Get payment details

## 🔐 Security Features

- **Password Security**: bcrypt hashing with salt rounds
- **JWT Authentication**: Secure token-based authentication
- **ObjectID Validation**: Proper MongoDB ObjectID validation
- **Input Validation**: Comprehensive input validation for all endpoints

## 🧪 Testing Strategy

The testing approach includes:

1. **Unit Tests**: Core functionality validation
   - ObjectID generation and validation
   - Password hashing security
   - JWT token operations
   - Business logic validation

2. **Integration Tests**: End-to-end workflow testing
   - Authentication flows
   - CRUD operations with ObjectID
   - Payment processing
   - API endpoint validation

3. **Performance Tests**: ObjectID operation benchmarks

## 📁 Key Data Models

All models use MongoDB ObjectID for consistency and performance:

```go
type User struct {
    ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Username string            `json:"username" bson:"username"`
    Email    string            `json:"email" bson:"email"`
    Password string            `json:"password" bson:"password"`
    Role     string            `json:"role" bson:"role"`
}

type Film struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title       string            `json:"title" bson:"title"`
    Description string            `json:"description" bson:"description"`
    Duration    int              `json:"duration" bson:"duration"`
    Genre       string           `json:"genre" bson:"genre"`
}

// Additional models: Jadwal, Tiket, Pembayaran
```

## 🚀 Development

### Building
```bash
# Build the application
go build -o bin/main main.go

# Run with Makefile
make build
make run
```

### Contributing
1. Create a feature branch
2. Write tests for new functionality
3. Ensure all tests pass
4. Update documentation as needed
5. Submit a pull request

## 📝 Documentation

Detailed documentation is available in the `docs/` folder:

- **TESTING_PLAN.md**: Complete testing strategy and methodology
- **TEST_REPORT.md**: Testing results and coverage analysis
- **CONTROLLER_TEST_REPORT.md**: Controller-specific test documentation
- **README_POSTMAN.md**: API usage examples and Postman collection

## 🔍 Monitoring & Debugging

- Use `go test -v` for detailed test output
- Check logs for ObjectID validation errors
- Monitor JWT token expiration and refresh
- Validate MongoDB ObjectID format in API requests

## 📈 Performance

- ObjectID operations are optimized for MongoDB performance
- JWT tokens include ObjectID-to-string conversion for compatibility
- Test benchmarks ensure efficient ObjectID handling

## 🤝 Support

For issues and questions:
1. Check the documentation in `docs/`
2. Review test cases for usage examples
3. Validate ObjectID format requirements
4. Ensure proper JWT authentication headers

---

**Note**: This project has undergone complete ObjectID migration and includes comprehensive testing coverage to ensure reliability and maintainability.
    "duration": 181
  }'
```

#### Get All Films
```bash
curl -X GET http://localhost:3000/api/films
```

#### Book Ticket
```bash
curl -X POST http://localhost:3000/api/tikets \
  -H "Content-Type: application/json" \
  -d '{
    "id": "tiket001",
    "jadwal_id": "jadwal001",
    "nama": "John Doe",
    "email": "john@example.com",
    "jumlah": 2,
    "user_id": "user001"
  }'
```

## Development

### Setup
1. Clone repository
2. Install dependencies: `go mod tidy`
3. Setup environment variables in `.env`
4. Run: `go run main.go`

### Generate Swagger Docs
```bash
swag init
```

### Environment Variables
```env
MONGOSTRING=mongodb+srv://...
MIDTRANS_SERVER_KEY=SB-Mid-server-...
MIDTRANS_CLIENT_KEY=SB-Mid-client-...
MIDTRANS_ENVIRONMENT=sandbox
```

## Database Schema

### Collections
- **films**: Film information
- **jadwals**: Movie schedules
- **tikets**: Ticket bookings
- **users**: User accounts
- **pembayarans**: Payment records

## Contributing
1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request
