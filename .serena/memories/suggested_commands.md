# Suggested Commands for BananaPhone Development

## Building

### Build for Windows (from macOS/Linux)
```bash
GOOS=windows GOARCH=amd64 go build ./pkg/BananaPhone
```

### Build specific example
```bash
GOOS=windows GOARCH=amd64 go build ./example/simplealloc
```

### Build mkdirectwinsyscall tool
```bash
go build ./cmd/mkdirectwinsyscall
```

## Code Generation

### Generate syscall wrappers
```bash
go generate ./example/hideexample/banana
```

### Run mkdirectwinsyscall manually
```bash
go run ./cmd/mkdirectwinsyscall -output zsyscall_windows.go -mode auto syscalls.go
```

## Testing

### Run tests (will fail on non-Windows)
```bash
go test ./...
```

### Run tests for Windows target
```bash
GOOS=windows go test ./...
```

**Note**: Tests will fail on macOS/Linux due to Windows-specific dependencies. This is expected.

## Formatting

### Format all Go code
```bash
gofmt -w .
```

### Check formatting without modifying
```bash
gofmt -l .
```

## Dependencies

### Download dependencies
```bash
go mod download
```

### Tidy dependencies
```bash
go mod tidy
```

### Vendor dependencies (recommended for stability)
```bash
go mod vendor
```

## Git Commands (macOS/Darwin)

### Standard git operations
```bash
git status
git add <files>
git commit -m "message"
git push
git pull
```

### View history
```bash
git log --oneline -10
```

## macOS/Darwin Utility Commands

### File operations
- `ls -la` - List files with details
- `find . -name "*.go"` - Find Go files
- `rg "pattern"` - Search with ripgrep (faster than grep)
- `cat file.go` - View file contents
- `head -n 20 file.go` - View first 20 lines

### Directory operations
- `pwd` - Print working directory
- `cd path` - Change directory
- `mkdir dirname` - Create directory

## Development Workflow

### Typical workflow for adding new syscall
1. Create or edit `syscalls.go` with `//dsys` annotations
2. Run `go generate` to create wrapper code
3. Build for Windows: `GOOS=windows go build`
4. Test on Windows system

### Example syscall annotation
```go
//dsys NtAllocateVirtualMemory(hProcess uintptr, lpAddress *uintptr, zerobits uintptr, dwSize *uint32, flAllocationType uint32, flProtect uint32) (err error) = ntdll.NtAllocateVirtualMemory
```

## Notes
- Always set `GOOS=windows` when building/testing
- Use `go generate` for code generation from annotated files
- No linting configuration found - use standard `gofmt`
- Vendor dependencies for production use (API not stable)
