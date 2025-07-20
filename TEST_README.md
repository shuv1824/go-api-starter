# Unit Tests Documentation

This document explains the unit tests that have been created for the DDD template project.

## Test Structure

The unit tests are organized following the same structure as the main application:

```
internal/
├── common/auth/jwt_test.go          # JWT authentication service tests
└── domains/user/infra/
    ├── service_test.go              # User service tests (business logic)
    └── repository_test.go           # User repository tests (database operations)
```

## Test Configuration

### Configuration File
Tests use a dedicated test configuration file: `config.test.yaml`

```yaml
mode: test
port: 8081
secret: test-secret-key
database:
  type: postgres
  host: localhost
  port: 5432
  username: postgres
  password: 123456
  dbname: gostarter_test
  sslmode: disable
```

### Prerequisites
1. **PostgreSQL Database**: Tests require a running PostgreSQL instance
2. **Test Database**: The tests will create and clean up temporary databases automatically
3. **Goose Migrations**: Tests use the existing migration files in `internal/migration/schema/`

## Test Features

### 1. JWT Authentication Tests (`internal/common/auth/jwt_test.go`)
- **Token Generation**: Validates JWT token creation with proper claims
- **Token Validation**: Tests token parsing and validation
- **Token Expiration**: Verifies token expiry functionality
- **Security**: Tests wrong signatures and invalid tokens

### 2. User Service Tests (`internal/domains/user/infra/service_test.go`)
- **User Registration**: Tests successful registration and duplicate email handling
- **User Login**: Tests authentication with valid/invalid credentials
- **Password Hashing**: Ensures passwords are properly hashed
- **Mock Repository**: Uses mocks to isolate business logic testing

### 3. User Repository Tests (`internal/domains/user/infra/repository_test.go`)
- **Database Operations**: Tests all CRUD operations
- **Real Database**: Uses PostgreSQL with goose migrations
- **Data Isolation**: Each test gets a unique temporary database
- **Automatic Cleanup**: Databases are cleaned up after tests

## Running Tests

### Prerequisites Setup
1. **Start PostgreSQL**:
   ```bash
   # Make sure PostgreSQL is running and accessible with the credentials in config.test.yaml
   sudo service postgresql start
   ```

2. **Create Test Database User** (if needed):
   ```sql
   CREATE USER postgres WITH PASSWORD '123456';
   ALTER USER postgres CREATEDB;
   ```

### Run All Tests
```bash
go test ./...
```

### Run Specific Test Files
```bash
# JWT tests
go test ./internal/common/auth/

# User service tests (with mocks)
go test ./internal/domains/user/infra/ -run TestService

# User repository tests (with real database)
go test ./internal/domains/user/infra/ -run TestUserRepository
```

### Run Tests with Verbose Output
```bash
go test -v ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

## Test Database Management

The repository tests automatically:

1. **Create Unique Databases**: Each test creates a unique database (e.g., `gostarter_test_a1b2c3d4`)
2. **Run Migrations**: Uses goose to apply schema migrations
3. **Execute Tests**: Runs all database operations
4. **Cleanup**: Drops the test database automatically

## Troubleshooting

### Database Connection Issues
If tests fail with database connection errors:

1. **Check PostgreSQL Status**:
   ```bash
   sudo service postgresql status
   ```

2. **Verify Credentials**: Ensure the database credentials in `config.test.yaml` are correct

3. **Check Database Permissions**: The test user needs CREATE DATABASE permissions

### Config File Issues
If tests skip due to missing config:

1. **Verify File Path**: Ensure `config.test.yaml` exists in the project root
2. **Check File Permissions**: Make sure the file is readable

### Migration Issues
If migrations fail:

1. **Check Migration Files**: Ensure files exist in `internal/migration/schema/`
2. **Verify Goose Syntax**: Check that migration files have proper `-- +goose Up/Down` comments

## Mock vs Real Database Tests

### Service Tests (Mocks)
- **Fast**: No database overhead
- **Isolated**: Tests only business logic
- **Predictable**: Full control over repository behavior

### Repository Tests (Real Database)
- **Accurate**: Tests actual database operations
- **Migration Testing**: Verifies schema compatibility
- **Integration**: Tests real database constraints and features

## Test Coverage

The tests cover:

- ✅ User registration and login flow
- ✅ Password hashing and verification
- ✅ JWT token generation and validation
- ✅ Database CRUD operations
- ✅ Error handling and edge cases
- ✅ Soft delete functionality
- ✅ Database migrations
- ✅ Configuration loading

## Contributing

When adding new tests:

1. **Follow Naming Conventions**: Use `Test<FunctionName>` format
2. **Use Table-Driven Tests**: For multiple test cases
3. **Mock External Dependencies**: Use mocks for service layer tests
4. **Clean Up Resources**: Ensure proper cleanup in database tests
5. **Document Complex Logic**: Add comments for complex test scenarios
