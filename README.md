# Go API Starter

A production-ready REST API template built with Gin, GORM, and supporting multiple database backends. This project follows Domain-Driven Design (DDD) principles and provides a solid foundation for building scalable Go web applications.

## Features

- **Fast HTTP Framework**: Built with [Gin](https://github.com/gin-gonic/gin) for high performance
- **Multiple Database Support**: PostgreSQL, MySQL, and SQLite through GORM
- **Clean Architecture**: Follows DDD principles with clear separation of concerns
- **Configuration Management**: YAML-based configuration with environment support
- **Middleware Support**: Built-in CORS and logging middleware
- **CLI Interface**: Powered by Cobra for command-line operations
- **Docker Ready**: Includes Dockerfile for containerization

## Project Structure

```
├── cmd/                  # Command line interface
│   └── root.go           # Root command and application bootstrap
├── internal/             # Private application code
│   ├── common/           # Shared utilities and middleware
│   │   ├── errors/       # Error handling utilities
│   │   └── middleware/   # HTTP middleware (CORS, logging)
│   └── config/           # Configuration management
├── pkg/                  # Public packages
│   └── database/         # Database abstraction layer
├── config.yaml           # Application configuration
├── Dockerfile            # Container configuration
├── Makefile              # Build automation
└── main.go               # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.23.6 or higher
- PostgreSQL, MySQL, or SQLite (depending on your choice)
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/shuv1824/go-api-starter.git
   cd go-api-starter
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Configure the application:

   ```bash
   cp config.yaml.example config.yaml  # if example exists
   # Edit config.yaml with your database settings
   ```

### Configuration

The application uses a YAML configuration file (`config.yaml`). Here's the default structure:

```yaml
mode: debug # Application mode: debug, test, release
port: 8080 # Server port

database:
  type: postgres # Database type: postgres, mysql, sqlite
  host: localhost
  port: 5432
  username: postgres
  password: 123456
  dbname: gostarter
  sslmode: disable
```

### Running the Application

#### Using Make

```bash
# Build and run
make

# Build only
make build
```

#### Using Go directly

```bash
# Run directly
go run main.go

# Build and run binary
go build -o bin/apiserver .
./bin/apiserver
```

### API Endpoints

The application includes a health check endpoint:

- `GET /ping` - Returns a simple pong response

## Database Support

The application supports multiple database backends through a factory pattern:

- **PostgreSQL** (default)
- **MySQL**
- **SQLite**

Simply change the `database.type` in your configuration file to switch between databases.

## Development

### Project Layout

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

- `cmd/` - Main applications for this project
- `internal/` - Private application and library code
- `pkg/` - Library code that's safe to use by external applications

### Adding New Features

1. **Domain Logic**: Add business logic in appropriate packages under `internal/`
2. **Database Models**: Define GORM models and repositories in domain packages
3. **HTTP Handlers**: Create REST endpoints in handler packages
4. **Middleware**: Add custom middleware in `internal/common/middleware/`

### Middleware

The application includes several built-in middleware:

- **CORS**: Cross-origin resource sharing support
- **Logging**: Request/response logging
- **Recovery**: Panic recovery middleware

## Docker Support

Build and run with Docker:

```bash
# Build Docker image
docker build -t go-api-starter .

# Run container
docker run -p 8080:8080 go-api-starter
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library for Go
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
