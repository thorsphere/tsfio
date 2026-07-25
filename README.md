# tsfio

[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/thorsphere/tsfio)](https://pkg.go.dev/mod/github.com/thorsphere/tsfio)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thorsphere/tsfio)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/thorsphere/tsfio)
![GitHub Top Language](https://img.shields.io/github/languages/top/thorsphere/tsfio)
[![CodeFactor](https://www.codefactor.io/repository/github/thorsphere/tsfio/badge)](https://www.codefactor.io/repository/github/thorsphere/tsfio)
![OSS Lifecycle](https://img.shields.io/osslifecycle/thorsphere/tsfio)

Go package providing a lightweight, safe API for file input/output operations. tsfio extends the standard library with utility functions for common file operations, test workflows, and string handling—built with security and simplicity in mind.

- **Zero-Config**: Use out-of-the-box with sensible defaults, no configuration required
- **Secure**: System-critical directories and files are blocked on Linux and Windows (see [inval_unix.go](https://github.com/thorsphere/tsfio/blob/main/inval_unix.go) and [inval_win.go](https://github.com/thorsphere/tsfio/blob/main/inval_win.go))
- **Tested**: Comprehensive unit tests with high code coverage
- **Minimal Dependencies**: Only requires the [Go Standard Library](https://pkg.go.dev/std) and [tserr](https://github.com/thorsphere/tserr)

## Safety & Defaults

File operations validate paths and return errors for system-critical locations. All file and directory operations use sensible defaults:

- Files are opened with read-write, append, and create flags (os.O_RDWR | os.O_APPEND | os.O_CREATE)
- File permissions: 0644
- Directory permissions: 0755
- All errors are returned in JSON format using [tserr](https://github.com/thorsphere/tserr)

## API

Import the package:

```go
import "github.com/thorsphere/tsfio"
```

### Type Definitions

Define file and directory paths:

```go
type Filename string  // Path to a regular file
type Directory string // Path to a directory
```

### Validation

All external functions validate input before operating. Validation functions are available:

```go
func CheckFile(f Filename) error // Validate file paths
func CheckDir(d Directory) error  // Validate directory paths
```

### File Operations

```go
func OpenFile(fn Filename) (*os.File, error)           // Open/create a file
func CloseFile(f *os.File) error                       // Close a file
func WriteStr(fn Filename, s string) error             // Append string to file
func WriteSingleStr(fn Filename, s string) error       // Write string (truncate if exists)
func TouchFile(fn Filename) error                      // Create/update file timestamp
func ReadFile(f Filename) ([]byte, error)              // Read file contents
func RemoveFile(f Filename) error                      // Delete a file
func ResetFile(fn Filename) error                      // Truncate file to empty
func ExistsFile(fn Filename) (bool, error)             // Check if file exists
func FileSize(fn Filename) (int64, error)              // Get file size in bytes
```

### Directory Operations

```go
func CreateDir(d Directory) error                      // Create directory (with parents)
func ExistsDir(dn Directory) (bool, error)             // Check if directory exists
func RemoveDir(d Directory) error                      // Remove directory and all its contents
```

### File Copy & Append

```go
type Copy struct {
	Src Filename // Source file
	Dst Filename // Destination file
}
func CopyFile(a *Copy) error                           // Copy file contents (dst must not exist)

type Append struct {
	FileA Filename // File to extend
	FileI Filename // File to append
}
func AppendFile(a *Append) error                       // Append one file to another
```

### Printable Characters

Filter non-printable characters from strings and runes:

```go
func Printable(a string) string                        // Remove non-printable runes
func IsPrintable(a []string) (bool, error)             // Check if strings are printable
func RuneToPrintable(r rune) string                    // Convert rune to printable
```

### Golden Files

Manage reference outputs for regression testing:

```go
type Testcase struct {
	Name string // Test case name
	Data string // Test data
}
func GoldenFilePath(name string) (Filename, error)     // Get golden file path
func CreateGoldenFile(tc *Testcase) error              // Create reference golden file
func EvalGoldenFile(tc *Testcase) error                // Compare against golden file
```

### Line Ending Normalization

Normalize line endings across platforms (CRLF/CR → LF):

```go
func NormNewlinesBytes(i []byte) ([]byte, error)       // Normalize byte slice
func NormNewlinesStr(i string) string                  // Normalize string
```

## Path Helpers

```go
func Sprintf[T Fio](f string, a ...any) T          		// Sprintf for file and directory paths
func Path(d Directory, f Filename) (Filename, error) 	// Join directory and filename
```

## Blocked Path Inspection

```go
func InvalDir() []Directory  							// List of blocked directories
func InvalFile() []Filename  							// List of blocked filenames
```

## Quick Start

```go
package main

import (
	"fmt"
	"os"

	"github.com/thorsphere/tsfio"
)

func main() {
	// Create temporary test files
	f1, _ := os.CreateTemp("", "foo")
	fn1 := tsfio.Filename(f1.Name())
	defer f1.Close()

	f2, _ := os.CreateTemp("", "foo")
	fn2 := tsfio.Filename(f2.Name())
	defer f2.Close()

	// Remove and check if file exists
	tsfio.RemoveFile(fn1)
	exists, _ := tsfio.ExistsFile(fn1)
	fmt.Println(exists) // false

	// Create (touch) a file
	tsfio.TouchFile(fn1)
	exists, _ = tsfio.ExistsFile(fn1)
	fmt.Println(exists) // true

	// Append strings to file
	tsfio.WriteStr(fn1, "foo")
	tsfio.WriteStr(fn1, "foo")
	data, _ := tsfio.ReadFile(fn1)
	fmt.Println(string(data)) // foofoo

	// Append one file to another
	tsfio.WriteStr(fn2, "foo")
	tsfio.AppendFile(&tsfio.Append{FileA: fn1, FileI: fn2})
	data, _ = tsfio.ReadFile(fn1)
	fmt.Println(string(data)) // foofoofoo

	// Get file size
	size, _ := tsfio.FileSize(fn1)
	fmt.Println(size) // 9

	// Overwrite file with single string
	tsfio.WriteSingleStr(fn1, "foo")
	data, _ = tsfio.ReadFile(fn1)
	fmt.Println(string(data)) // foo

	// Reset (truncate) file
	tsfio.ResetFile(fn1)
	data, _ = tsfio.ReadFile(fn1)
	fmt.Println(string(data)) // (empty)
}
```

[Try on Go Playground](https://go.dev/play/p/NHvbqw1L3J1)

## Documentation & Resources

- [Go Package Documentation](https://pkg.go.dev/github.com/thorsphere/tsfio) — Complete API reference
- [Open Source Insights](https://deps.dev/go/github.com%2Fthorsphere%2Ftsfio) — Dependency analysis

## ⚖️ License & Commercial Usage

Copyright (c) 2023-2026 thorsphere. All rights reserved.

This project is licensed under the **Functional Source License v1.1 (FSL-1.1-ALv2)**. 

* The use, modification, and redistribution of this Go package is completely free for private, educational, non-commercial, and internal purposes. 
* If you are a company or institution looking to use this package in a commercial product, service, or business environment, you must secure a commercial license.
* Each version of this software automatically converts to the fully open-source Apache License, Version 2.0 on the second anniversary of its release.

For full details, please see the [LICENSE](LICENSE.md) file.

### 💼 Commercial Licensing & Inquiries

To purchase a commercial license or discuss support options, please reach out directly:

* 📩 **Contact:** business at thorsphere dot com
* 💬 **Response Time:** Usually within a couple of business days.

*Please include your company name and a brief overview of your use case so I can provide the right licensing details.*
