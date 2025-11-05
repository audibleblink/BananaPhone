# BananaPhone Project Overview

## Purpose
BananaPhone is a pure-Go implementation for performing direct Windows syscalls, bypassing standard API calls to evade detection by AV/EDR solutions. It's similar to Hell's Gate but implemented in Go.

**Key Concept**: By directly invoking syscalls with resolved system IDs (sysids), BananaPhone avoids API hooking mechanisms used by security monitoring tools like API Monitor.

## Target Platform
- **Windows only** - This library is specifically designed for Windows and will not work on other platforms
- Development can be done on macOS/Linux, but builds must target Windows (GOOS=windows)

## Core Features
1. **Direct Syscall Execution**: Call Windows kernel functions directly using resolved sysids
2. **Multiple Resolution Modes**:
   - Memory mode: Parse PEB in-memory to find ntdll and resolve exports
   - Disk mode: Load ntdll.dll from disk to resolve exports
   - Halo's Gate mode: Deduce syscalls from non-hooked neighboring functions
   - Auto mode: Try memory → Halo's Gate → disk (recommended)
3. **Low-level Functions**:
   - `GetPEB()`: Get Process Environment Block without API calls
   - `GetNtdllStart()`: Get ntdll start address from memory
   - `WriteMemory()`: Write bytes to memory address
   - `Syscall()`: Execute syscalls with resolved sysids
4. **Code Generation Tool**: `mkdirectwinsyscall` - generates Go code for direct syscalls from annotated function signatures

## Tech Stack
- **Language**: Go 1.15+
- **Key Dependencies**:
  - `github.com/Binject/debug` - PE parsing and manipulation
  - `github.com/awgh/rawreader` - Raw memory reading
  - `golang.org/x/sys/windows` - Windows syscall support
- **Assembly**: x64 assembly (asm_x64.s) for low-level operations

## Repository Structure
```
BananaPhone/
├── pkg/BananaPhone/          # Core library
│   ├── bananaphone.go        # Main API and resolution logic
│   ├── functions.go          # Syscall wrapper functions
│   ├── internal.go           # Internal helpers
│   ├── ldr.go                # PE loader functions
│   └── asm_x64.s             # x64 assembly routines
├── cmd/mkdirectwinsyscall/   # Code generation tool
├── example/                  # Usage examples
│   ├── simplealloc/          # Basic memory allocation
│   ├── hideexample/          # Comparison with/without BananaPhone
│   ├── mkwinsyscall/         # Code generation example
│   └── ...                   # Other examples
└── .github/                  # GitHub templates
```

## Important Notes
- **API Stability**: API is not yet stable - vendor dependencies properly
- **Windows-only**: Will not work on macOS/Linux at runtime
- **Security Research**: This is a security research/red team tool
