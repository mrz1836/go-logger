# CLAUDE.md - go-logger

## 🎯 Project Overview
**go-logger** is a high-performance Go logging library that extends the native log package with automatic Log Entries/Rapid7 integration, queuing, and file tracking.

### Core Architecture
- **Auto-switching implementation**: Detects `LOG_ENTRIES_TOKEN` env var to switch between local logging and Rapid7 integration
- **Message queuing**: Thread-safe TCP connection with retry logic and exponential backoff
- **File tracking**: Automatic file/line/method detection using `runtime.Caller()`
- **Zero dependencies**: Only testify for tests, no runtime dependencies

## 📁 Key Files
- `logger.go:17-27` - Main Logger interface and global functions
- `log_entries.go:38-46` - LogClient with TCP connection and queue management
- `gorm.go:12-22` - GORM compatibility interface for database logging
- `parameters.go:7-11` - Parameter struct for structured logging
- `config.go:8-14` - Constants for Rapid7 endpoints and retry settings

## 🔧 Build Commands
```bash
magex test           # Run all tests (coverage threshold: 65%)
magex test:race      # Run with race detector
magex lint           # Run golangci-lint
magex format:fix     # Format code with gofumpt
magex bench          # Run benchmarks
magex deps:update    # Update all dependencies
```

## 🌟 Key Implementation Patterns

### Environment Configuration
```go
LOG_ENTRIES_TOKEN=your-token        # Enables Rapid7 integration
LOG_ENTRIES_ENDPOINT=custom-url     # Optional custom endpoint
LOG_ENTRIES_PORT=10000              # Optional custom port
```

### Logger Interface Switching
- `init()` function automatically detects env vars and switches implementations
- `SetImplementation(Logger)` allows runtime switching
- Falls back to standard log package if no token provided

### Message Queue & Retry Logic
- `LogClient.ProcessQueue()` runs in goroutine for async processing
- Exponential backoff on connection failures (100ms → 2min max)
- `PushFront()` for message reordering on failures

### Structured Logging
```go
logger.Data(2, logger.INFO, "message",
    logger.MakeParameter("key", value))
```

## 🧪 Testing Approach
- Comprehensive unit tests with testify
- Fuzz testing for all major functions
- TCP connection mocking for Log Entries client
- GORM interface testing with mock contexts

## ⚠️ Important Notes
- Always use `magex test` before committing
- File tracking adds 2 stack levels to most logger calls
- Log Entries client requires network connectivity for tests
- GORM logger uses 5s threshold for slow query warnings

## 🔍 Common Code Locations
- Global logger functions: `logger.go:154-241`
- TCP connection management: `log_entries.go:66-99`
- Queue processing loop: `log_entries.go:102-122`
- GORM trace function: `gorm.go:106-142`
- File/line utilities: `logger.go:115-152`
