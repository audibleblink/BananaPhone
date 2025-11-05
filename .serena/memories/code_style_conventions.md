# Code Style and Conventions

## General Go Style
- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (standard Go formatter)
- No specific linter configuration found in repository

## Naming Conventions

### Types and Structs
- **PascalCase** for exported types: `BananaPhone`, `PhoneMode`, `BananaProcedure`
- **camelCase** for unexported types: `bpSyscall`
- Descriptive names that indicate purpose

### Constants
- **PascalCase** for exported constants with descriptive names
- Group related constants using `iota`:
```go
const (
    MemoryBananaPhoneMode PhoneMode = iota
    DiskBananaPhoneMode
    AutoBananaPhoneMode
    HalosGateBananaPhoneMode
)
```

### Functions and Methods
- **PascalCase** for exported functions: `NewBananaPhone`, `GetPEB`, `Syscall`
- **camelCase** for unexported functions: `getSysID`, `bpSyscall`
- Constructor pattern: `New<TypeName>` for constructors
- Pointer receivers for methods that modify state: `(*BananaPhone).GetSysID`
- Value receivers for methods that don't modify: `(BananaProcedure).Addr`

### Variables
- **camelCase** for local variables: `phandle`, `baseA`, `regionsize`
- Short, descriptive names in limited scopes
- More descriptive names for package-level variables

## Comments and Documentation

### Package Comments
- Package-level comment at top of main file
```go
package bananaphone
```

### Type Comments
- Comment above type definition explaining purpose
- Use `//TypeName` format:
```go
//PhoneMode determines the way a bananaphone will resolve sysids
type PhoneMode int
```

### Function Comments
- Comment above function with description
- Use `//FunctionName` format:
```go
//NewBananaPhone creates a new instance of a bananaphone with behaviour as defined by the input value. Use AutoBananaPhoneMode if you're not sure.
func NewBananaPhone(t PhoneMode) (*BananaPhone, error) {
```

### Inline Comments
- Used sparingly for complex logic
- Explain "why" not "what"

## Error Handling
- Return errors as last return value: `func NewBananaPhone(...) (*BananaPhone, error)`
- Use `fmt.Errorf` for error messages
- Custom error types for specific cases: `MayBeHookedError`
- Some functions panic for unrecoverable errors (e.g., `WriteMemory`)

## Code Organization
- Group related functionality in separate files:
  - `bananaphone.go` - Core API
  - `functions.go` - Syscall wrappers
  - `internal.go` - Internal helpers
  - `ldr.go` - PE loading
  - `asm_x64.s` - Assembly routines

## Assembly Code Style (asm_x64.s)
- Go assembly syntax (Plan 9 style)
- Function signature in comment above implementation
- Clear register usage
- Example:
```asm
//func GetPEB() uintptr
TEXT ·GetPEB(SB), $0-8
     MOVQ 	0x60(GS), AX
     MOVQ	AX, ret+0(FP)
     RET
```

## Import Organization
- Standard library imports first
- Third-party imports after
- Blank line between groups

## No Type Hints/Generics
- Project uses Go 1.15 (pre-generics)
- Explicit type declarations where needed
- Use `interface{}` for generic types (though rare in this codebase)
