# Finance Tracking API - Test Suite

Comprehensive test suite untuk menguji endpoint, arsitektur kode, alur business logic, dan clean code practices.

## 📁 Struktur Test

```
test/
├── README.md                           # Dokumentasi test suite
├── run_tests.ps1                       # Script untuk menjalankan tests
│
├── helpers/                            # Test utilities dan helper functions
│   ├── database_helper.go             # Setup test database (SQLite in-memory)
│   └── http_helper.go                 # HTTP request/response helpers
│
├── repositories/                       # Unit tests untuk data access layer
│   ├── income_repository_test.go      # Test CRUD operations untuk income
│   └── category_repository_test.go    # Test CRUD operations untuk categories
│
├── services/                           # Unit tests untuk business logic layer
│   └── income_service_test.go         # Test auto-allocation logic, validations
│
├── handlers/                           # Unit tests untuk API endpoints
│   └── category_handler_test.go       # Test HTTP handlers dan response formats
│
├── integration/                        # Integration tests (end-to-end)
│   └── api_integration_test.go        # Test complete workflows dari API
│
└── architecture/                       # Architecture dan code quality tests
    └── clean_architecture_test.go     # Test layer separation, naming conventions, DI
```

## 🎯 Jenis Tests

### 1. **Repository Tests** (`test/repositories/`)
Menguji data access layer dan database operations:
- ✅ CRUD operations (Create, Read, Update, Delete)
- ✅ Query dengan filters (by month/year, by type, etc.)
- ✅ Database constraints dan validations
- ✅ Error handling untuk record not found

**Coverage:**
- `income_repository_test.go`: 9 test cases
- `category_repository_test.go`: 8 test cases

### 2. **Service Tests** (`test/services/`)
Menguji business logic dan core features:
- ✅ Auto budget allocation (proportional distribution)
- ✅ Budget calculation dan validations
- ✅ Transaction handling
- ✅ Edge cases (no active categories, zero budgets, etc.)

**Coverage:**
- `income_service_test.go`: 4 test cases termasuk auto-allocation logic

### 3. **Handler Tests** (`test/handlers/`)
Menguji API endpoints dan HTTP layer:
- ✅ Request validation
- ✅ Response format (success/error)
- ✅ HTTP status codes
- ✅ JSON serialization/deserialization
- ✅ Parameter parsing (path params, query params)

**Coverage:**
- `category_handler_test.go`: 6 test cases untuk semua category endpoints

### 4. **Integration Tests** (`test/integration/`)
Menguji complete workflows end-to-end:
- ✅ Complete user journey (create category → income → expense)
- ✅ Auto-allocation flow
- ✅ Budget tracking accuracy
- ✅ Concurrent request handling
- ✅ Data consistency across layers

**Coverage:**
- `api_integration_test.go`: 2 comprehensive integration scenarios

### 5. **Architecture Tests** (`test/architecture/`)
Menguji code quality dan clean architecture principles:
- ✅ **Layer Separation**: Models, Repositories, Services, Handlers tidak cross-import
- ✅ **Naming Conventions**: Files follow `*_repository.go`, `*_service.go` patterns
- ✅ **Interface Definitions**: Repositories dan Services define interfaces
- ✅ **Error Handling**: Proper error returns
- ✅ **Dependency Injection**: Constructor-based DI pattern

**Coverage:**
- `clean_architecture_test.go`: 5 architectural validation tests

## 🚀 Cara Menjalankan Tests

### Menggunakan PowerShell Script (Recommended)

```powershell
# Run all tests
.\test\run_tests.ps1

# Run specific test type
.\test\run_tests.ps1 -TestType unit          # Repository + Service tests
.\test\run_tests.ps1 -TestType integration   # Integration tests only
.\test\run_tests.ps1 -TestType handlers      # Handler tests only
.\test\run_tests.ps1 -TestType architecture  # Architecture tests only

# Run with coverage report
.\test\run_tests.ps1 -TestType coverage      # Generates coverage.html
```

### Menggunakan Go Command Langsung

```bash
# Run all tests
go test ./test/... -v

# Run specific package
go test ./test/repositories/... -v
go test ./test/services/... -v
go test ./test/handlers/... -v
go test ./test/integration/... -v
go test ./test/architecture/... -v

# Run with coverage
go test ./test/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run specific test function
go test ./test/repositories/... -run TestIncomeRepository_Create -v
```

## 📊 Test Output Example

```
=== RUN   TestIncomeRepository_Create
--- PASS: TestIncomeRepository_Create (0.01s)

=== RUN   TestIncomeService_CreateIncome_WithAutoAllocation
--- PASS: TestIncomeService_CreateIncome_WithAutoAllocation (0.02s)

=== RUN   TestCategoryHandler_CreateCategory
=== RUN   TestCategoryHandler_CreateCategory/Valid_category_creation
=== RUN   TestCategoryHandler_CreateCategory/Missing_required_fields
--- PASS: TestCategoryHandler_CreateCategory (0.03s)

=== RUN   TestEndToEnd_CompleteWorkflow
=== RUN   TestEndToEnd_CompleteWorkflow/Create_Categories
=== RUN   TestEndToEnd_CompleteWorkflow/Get_All_Categories
=== RUN   TestEndToEnd_CompleteWorkflow/Create_Income_with_Auto-Allocation
--- PASS: TestEndToEnd_CompleteWorkflow (0.05s)

PASS
coverage: 85.3% of statements
```

## 🧪 Test Database

Tests menggunakan **SQLite in-memory database** untuk:
- ✅ Fast execution (no disk I/O)
- ✅ Isolated environment (setiap test independen)
- ✅ No cleanup required (auto-destroyed after test)
- ✅ Same schema as PostgreSQL production

## ✅ What is Tested

### Clean Architecture Principles
- [x] Layer separation (no circular dependencies)
- [x] Dependency injection pattern
- [x] Interface-based programming
- [x] Proper error handling
- [x] Naming conventions

### API Endpoints
- [x] `POST /api/v1/categories` - Create category
- [x] `GET /api/v1/categories` - List all categories
- [x] `GET /api/v1/categories/:id` - Get category by ID
- [x] `PUT /api/v1/categories/:id` - Update category
- [x] `DELETE /api/v1/categories/:id` - Delete category
- [x] `POST /api/v1/incomes` - Create income (auto-allocation)
- [x] `GET /api/v1/incomes` - List incomes
- [x] `POST /api/v1/expenses` - Create expense (auto-deduction)
- [x] `GET /api/v1/budgets` - Get category budgets
- [x] `GET /api/v1/budgets/summary` - Get budget summary

### Business Logic
- [x] Auto budget allocation (proportional distribution)
- [x] Budget deduction when expense created
- [x] Monthly budget tracking
- [x] Category priority handling
- [x] Active/inactive category filtering

### Edge Cases
- [x] Creating income without active categories
- [x] Expense amount exceeds budget
- [x] Invalid UUID formats
- [x] Missing required fields
- [x] Concurrent requests

## 📝 Best Practices Covered

1. **Unit Testing**: Each layer tested independently
2. **Integration Testing**: End-to-end scenarios
3. **Test Isolation**: No test depends on another
4. **Meaningful Assertions**: Clear error messages
5. **Test Data Management**: Helper functions for seeding
6. **Coverage**: Core business logic has high coverage

## 🔧 Dependencies

```bash
# Test dependencies
go get github.com/stretchr/testify/assert  # Assertion library
go get gorm.io/driver/sqlite               # In-memory database
```

## 📈 Coverage Goals

- **Repositories**: > 80% coverage
- **Services**: > 85% coverage (critical business logic)
- **Handlers**: > 75% coverage
- **Overall**: > 80% coverage

## 🐛 Debugging Tests

```bash
# Run with verbose output
go test ./test/... -v

# Run specific test with trace
go test ./test/repositories/... -run TestIncomeRepository_Create -v -trace trace.out

# Check for race conditions
go test ./test/... -race

# Run with timeout
go test ./test/... -timeout 30s
```

## 📚 Test Documentation

Setiap test function didokumentasikan dengan:
- **Setup**: Persiapan test environment
- **Execute**: Menjalankan function yang ditest
- **Assert**: Verifikasi hasil yang diharapkan

Example:
```go
func TestIncomeRepository_Create(t *testing.T) {
    // Setup - persiapan database dan test data
    db := helpers.SetupTestDB(t)
    defer helpers.CleanupTestDB(t, db)
    
    // Execute - jalankan function
    err := repo.Create(income)
    
    // Assert - verifikasi hasil
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, income.ID)
}
```

## 🎓 Learning Resources

- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Testing Package](https://pkg.go.dev/testing)
- [GORM Testing Guide](https://gorm.io/docs/testing.html)
- [Clean Architecture Testing](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 💡 Tips

1. **Run tests frequently** during development
2. **Write tests first** (TDD approach) untuk critical features
3. **Keep tests simple** dan focused pada satu aspek
4. **Use table-driven tests** untuk multiple scenarios
5. **Mock external dependencies** (database, API calls)

---

**Status**: ✅ All test suites implemented and passing
**Last Updated**: February 17, 2026
