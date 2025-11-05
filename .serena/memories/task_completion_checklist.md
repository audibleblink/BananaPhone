# Task Completion Checklist

When completing a task in the BananaPhone project, follow these steps:

## 1. Code Quality

### Formatting
- [ ] Run `gofmt -w .` to format all modified Go files
- [ ] Verify formatting with `gofmt -l .` (should return nothing)

### Code Review
- [ ] Ensure code follows project conventions (see code_style_conventions.md)
- [ ] Check that exported functions/types have proper comments
- [ ] Verify error handling is consistent with project patterns
- [ ] Ensure Windows-specific code is properly guarded

## 2. Building

### Verify Builds
- [ ] Build for Windows target: `GOOS=windows GOARCH=amd64 go build ./pkg/BananaPhone`
- [ ] If adding examples, build them: `GOOS=windows go build ./example/<name>`
- [ ] If modifying mkdirectwinsyscall, build it: `go build ./cmd/mkdirectwinsyscall`

### Code Generation
- [ ] If syscall annotations changed, run `go generate` on affected packages
- [ ] Verify generated files are correct and committed

## 3. Dependencies

### Module Management
- [ ] Run `go mod tidy` to clean up dependencies
- [ ] Verify `go.mod` and `go.sum` are updated if dependencies changed
- [ ] Consider running `go mod vendor` for stability

## 4. Testing

### Test Execution
- [ ] Run `GOOS=windows go test ./...` (will have limited success on non-Windows)
- [ ] If possible, test on actual Windows system
- [ ] Verify examples still compile: `GOOS=windows go build ./example/...`

**Note**: Full testing requires Windows environment due to platform-specific nature

## 5. Documentation

### Code Documentation
- [ ] Add/update function comments for new/modified exported functions
- [ ] Update README.md if adding new features or changing API
- [ ] Update cmd/mkdirectwinsyscall/readme.md if changing tool behavior
- [ ] Add examples if introducing new functionality

### Memory Files
- [ ] Update memory files if project structure or conventions change

## 6. Git

### Commit Preparation
- [ ] Review changes: `git diff`
- [ ] Stage relevant files: `git add <files>`
- [ ] Check status: `git status`
- [ ] Commit with descriptive message: `git commit -m "description"`

### Commit Message Guidelines
- Use present tense: "Add feature" not "Added feature"
- Be specific: "Add Halo's Gate mode support" not "Update code"
- Reference issues if applicable: "Fix #123: ..."

## 7. Platform Considerations

### Windows-Only Code
- [ ] Verify build tags if needed: `//go:build windows`
- [ ] Ensure assembly code is x64 compatible
- [ ] Check that syscall IDs are correct for target Windows version

### Cross-Platform Development
- [ ] Remember this code only runs on Windows
- [ ] Development/building can be done on macOS/Linux with GOOS=windows
- [ ] Testing requires actual Windows environment

## 8. API Stability Warning

### Breaking Changes
- [ ] Note that API is not yet stable
- [ ] Document breaking changes clearly
- [ ] Consider impact on existing users
- [ ] Update examples if API changes affect them

## Quick Checklist for Small Changes

For minor changes (typos, comments, small fixes):
1. Format: `gofmt -w .`
2. Build: `GOOS=windows go build ./pkg/BananaPhone`
3. Tidy: `go mod tidy`
4. Commit: `git add . && git commit -m "description"`

## Quick Checklist for Major Changes

For significant changes (new features, refactoring):
1. All items from "Quick Checklist for Small Changes"
2. Run code generation if needed: `go generate ./...`
3. Build all examples: `GOOS=windows go build ./example/...`
4. Update documentation (README, comments, memory files)
5. Test on Windows if possible
6. Review all changes carefully before committing
