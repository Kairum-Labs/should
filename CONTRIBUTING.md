# Contributing to Should

Thank you for your interest in contributing to Should! This guide covers everything you need to know to contribute effectively.

## Quick Start

### Prerequisites

- Go 1.22 or later
- Git

### Setup

1. Fork and clone the repository:

   ```bash
   git clone https://github.com/YOUR_USERNAME/should.git
   cd should
   ```

2. Run tests to ensure everything works:
   ```bash
   go test ./...
   ```

## Project Structure

```
should/
├── assert/
│   ├── assertions.go     # Main assertion methods
│   ├── types.go         # Type definitions
│   ├── utils.go         # Error formatting utilities
│   └── *_test.go        # Test files
├── README.md
└── go.mod
```

## How to Contribute

### Types of Contributions

- **Bug Fixes**: Fix existing issues
- **New Assertions**: Add new assertion methods
- **Performance**: Optimize existing code
- **Documentation**: Improve examples and docs
- **Tests**: Add test cases or improve coverage

### Before You Start

1. Check existing issues for similar work
2. For major changes, create an issue first to discuss
3. Keep changes focused (one feature per PR)

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Follow Go conventions and `gofmt`
- Add tests for new functionality
- Ensure all tests pass: `go test ./...`

## Linting (golangci-lint)

- CI: we use `golangci-lint` version `v2.3.0` (via `golangci/golangci-lint-action@v8`), configured by `.golangci.yml`.

### Local installation

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.0
```

### How to run

```bash
golangci-lint run ./...
```

Tips:

- Run in a specific directory: `golangci-lint run ./assert/...`
- Run on a specific file: `golangci-lint run assert/utils.go`

### 3. Commit

Use clear commit messages:

```bash
git commit -m "feat: add BeInRange assertion"
git commit -m "fix: handle nil pointers in BeEmpty"
```

**Commit Types:**

- `feat`: New feature
- `fix`: Bug fix
- `perf`: Performance improvement
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `style`: Code style changes

### 4. Submit PR

- Push to your fork
- Create a pull request with clear description
- Ensure CI passes

## Code Standards

### Assertion Methods

- Use "Be" prefix: `BeEqual`, `BeTrue`, `BeEmpty`
- Support custom messages: `opts ...Option`
- Always call `t.Helper()` for proper stack traces

```go
func BeEqual[T any](t testing.TB, actual T, expected T, opts ...Option) {
    t.Helper()
    if !reflect.DeepEqual(a.value, expected) {
        fail(t, formatError(a.value, expected))
    }
}
```

### Testing

- Test both success and failure cases
- Use descriptive test names: `TestBeEqual_WithIdenticalValues_Succeeds`
- Maintain >90% test coverage
- Include benchmarks for performance-critical code

### Error Messages

- Provide detailed, helpful error messages
- Include context and suggestions when possible
- Use consistent formatting

## Running Tests

```bash
# All tests
go test ./...

# With coverage and race detector (like CI)
go test -v -race -cover ./...

# Specific test
go test -run TestBeEqual ./...

# Format code
go fmt ./...
```

## Getting Help

- **Issues**: Use GitHub issues for bugs and feature requests
- **Questions**: Start a GitHub Discussion

## Recognition

All contributors will be recognized in release notes. Thank you for helping make Should better!
