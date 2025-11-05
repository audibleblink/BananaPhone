# BananaPhone Architecture and Design

## Core Concepts

### Direct Syscalls
BananaPhone bypasses Windows API hooking by directly invoking syscalls using resolved system IDs (sysids). This evades detection by security tools that hook standard API functions.

### System ID Resolution
The library resolves syscall IDs through multiple methods:
1. **Memory Mode**: Parse Process Environment Block (PEB) in-memory
2. **Disk Mode**: Load ntdll.dll from disk
3. **Halo's Gate Mode**: Deduce syscalls from neighboring non-hooked functions
4. **Auto Mode**: Try memory → Halo's Gate → disk (fallback chain)

## Key Components

### 1. BananaPhone Struct (`bananaphone.go`)
Core type that manages syscall resolution:
```go
type BananaPhone struct {
    banana *pe.File    // Parsed PE file (ntdll)
    mode   PhoneMode   // Resolution mode
    memloc uintptr     // Memory location of ntdll
}
```

**Key Methods**:
- `NewBananaPhone(mode)` - Constructor with mode selection
- `GetSysID(name)` - Resolve syscall ID by function name
- `NewProc(name)` - Get procedure by name
- `GetFuncPtr(name)` - Get function pointer

### 2. Syscall Execution (`functions.go`)
Low-level syscall invocation:
```go
func Syscall(callid uint16, argh ...uintptr) (errcode uint32, err error)
```
- Takes syscall ID and variable arguments
- Returns error code and optional error
- Implemented in assembly for direct syscall

### 3. Assembly Routines (`asm_x64.s`)
Critical low-level operations in x64 assembly:

**GetPEB()**: Retrieve Process Environment Block
```asm
MOVQ 0x60(GS), AX  // PEB is at GS:0x60
```

**GetNtdllStart()**: Find ntdll base address
- Walks PEB → LDR → InMemoryOrderModuleList
- Returns start address and size

**bpSyscall()**: Execute syscall with resolved ID
- Sets up syscall number and parameters
- Performs direct syscall instruction

### 4. PE Parsing (`ldr.go`, `internal.go`)
Uses `github.com/Binject/debug/pe` to:
- Parse ntdll.dll PE structure
- Extract export table
- Resolve function addresses
- Read syscall IDs from function prologue

### 5. Code Generation Tool (`cmd/mkdirectwinsyscall/`)
Generates Go wrapper code from annotated function signatures:

**Input** (`syscalls.go`):
```go
//dsys NtAllocateVirtualMemory(hProcess uintptr, ...) (err error) = ntdll.NtAllocateVirtualMemory
```

**Output** (`zsyscall_windows.go`):
- Generated wrapper function
- Syscall ID resolution
- Parameter marshaling
- Error handling

**Flags**:
- `-mode`: Resolution mode (auto/memory/disk/raw)
- `-noglobal`: Don't use global var for resolution
- `-trace`: Debug output for syscalls
- `-output`: Output filename

## Design Patterns

### 1. Mode Pattern
`PhoneMode` enum determines resolution strategy:
```go
type PhoneMode int
const (
    MemoryBananaPhoneMode PhoneMode = iota
    DiskBananaPhoneMode
    AutoBananaPhoneMode
    HalosGateBananaPhoneMode
)
```

### 2. Lazy Initialization
Syscall IDs resolved on first use, cached for subsequent calls (when using globals)

### 3. Error Handling
- Custom error types: `MayBeHookedError`
- Standard Go error returns
- Panic for unrecoverable errors

### 4. Unsafe Operations
Heavy use of `unsafe` package for:
- Memory manipulation
- Pointer arithmetic
- Direct memory writes

## Hook Detection

### Halo's Gate Implementation
When a function is hooked, BananaPhone can:
1. Detect the hook (unexpected bytes in function prologue)
2. Search neighboring functions for non-hooked versions
3. Deduce the syscall ID from neighboring function IDs
4. Calculate target syscall ID based on position

### Hook Check
```go
type MayBeHookedError struct {
    // Error indicating potential hook detected
}
```

## Memory Safety Considerations

### Unsafe Operations
- Direct memory access via `unsafe.Pointer`
- No bounds checking on memory writes
- Assumes valid memory addresses

### Platform Assumptions
- x64 architecture only
- Windows-specific memory layout
- PEB structure at GS:0x60

## Integration Points

### With Go Standard Library
- Minimal use of `syscall` package
- Uses `unsafe` for low-level operations
- Standard error handling patterns

### With Windows
- Direct interaction with Windows kernel
- No API calls for critical operations
- Reads ntdll.dll structure directly

## Performance Characteristics

### Resolution Overhead
- **Memory Mode**: Fast, in-memory parsing
- **Disk Mode**: Slower, file I/O required
- **Halo's Gate**: Medium, requires searching
- **Auto Mode**: Variable, depends on fallback

### Caching
- Global mode: Resolve once, cache sysid
- No-global mode: Resolve every call (slower)

## Security Considerations

### Evasion Techniques
- Bypasses API hooks
- Minimal API calls
- Direct syscall execution
- Multiple resolution fallbacks

### Detection Vectors
- File handle to ntdll.dll (disk mode)
- Memory scanning for syscall patterns
- Behavioral analysis of direct syscalls

## Extension Points

### Adding New Syscalls
1. Add `//dsys` annotation in `syscalls.go`
2. Run `go generate`
3. Generated wrapper handles resolution

### Custom Resolution
- Implement custom `PhoneMode`
- Extend `getSysID` logic
- Add new resolution strategy

## Dependencies

### Critical Dependencies
- `github.com/Binject/debug` - PE parsing
- `github.com/awgh/rawreader` - Raw memory reading
- `golang.org/x/sys/windows` - Windows types (build time only)

### Why These Dependencies
- **Binject/debug**: Extends stdlib `pe` with advanced parsing
- **rawreader**: Safe raw memory reading
- **x/sys/windows**: Type definitions and constants
