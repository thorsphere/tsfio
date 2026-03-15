# tsfio

[![Go Report Card](https://goreportcard.com/badge/github.com/thorsphere/tsfio)](https://goreportcard.com/report/github.com/thorsphere/tsfio)
[![CodeFactor](https://www.codefactor.io/repository/github/thorsphere/tsfio/badge)](https://www.codefactor.io/repository/github/thorsphere/tsfio)
![OSS Lifecycle](https://img.shields.io/osslifecycle/thorsphere/tsfio)

[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/thorsphere/tsfio)](https://pkg.go.dev/mod/github.com/thorsphere/tsfio)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thorsphere/tsfio)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/thorsphere/tsfio)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/thorsphere/tsfio)
![GitHub last commit](https://img.shields.io/github/last-commit/thorsphere/tsfio)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/thorsphere/tsfio)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/thorsphere/tsfio)
![GitHub Top Language](https://img.shields.io/github/languages/top/thorsphere/tsfio)
![GitHub](https://img.shields.io/github/license/thorsphere/tsfio)

[Go](https://go.dev/) package with a [simple](https://en.wikipedia.org/wiki/KISS_principle) API for file input output. It is a supplement to the standard library and supplies additional functions for file input output and string operations, e.g., appending one file to another file.

- **Simple**: Without configuration, just function calls, and default flags are used
- **Resilient**: File input output on Linux and Windows system directories or files are blocked (see [inval_unix.go](https://github.com/thorstenrie/tsfio/blob/main/inval_unix.go) and [inval_win.go](https://github.com/thorstenrie/tsfio/blob/main/inval_win.go))
- **Tested**: Unit tests with a high code coverage
- **Dependencies**: Only depends on the [Go Standard Library](https://pkg.go.dev/std) and [tserr](https://github.com/thorstenrie/tserr)

## Defaults

All file input output operations on Linux and Windows system directories or
files are blocked (see [inval_unix.go](https://github.com/thorsphere/tsfio/blob/main/inval_unix.go) and [inval_win.go](https://github.com/thorsphere/tsfio/blob/main/inval_win.go)) and an error is returned.
All operations expect a directory or a regular file, return an error otherwise.
Default flags and file mode is used when opening files, creating files or directories
and when writing to files (with exceptions documented in the function descriptions)

- Files are opened read-write (os.O_RDWR).
- Data is appended when writing to file (os.O_APPEND).
- A file is created if it does not exist (os.O_CREATE).
- File mode and permission bits are 0644.
- Directory mode and permissions bits are 0755.

If an API call is not successful, a [tserr](https://github.com/thorsphere/tserr) error in JSON format is returned.

## Usage

In the Go app, the package is imported with

```go
import "github.com/thorsphere/tsfio"
```

A Filename is the name of a regular file and may contain its path. A Directory is the name of a directory and may contain its path

```go
type Filename string
type Directory string
```

CheckFile performs checks on Filename f and CheckDir performs checks on Directory d

```go
func CheckFile(f Filename) error
func CheckDir(d Directory) error
```

All external functions contain a CheckFile or CheckDir call at the beginning.

```go
func OpenFile(fn Filename) (*os.File, error)
func CloseFile(f *os.File) error
func WriteStr(fn Filename, s string) error
func WriteSingleStr(fn Filename, s string) error
func TouchFile(fn Filename) error
func ReadFile(f Filename) ([]byte, error)
func AppendFile(a *Append) error
func ExistsFile(fn Filename) (bool, error)
func RemoveFile(f Filename) error
func ResetFile(fn Filename) error
func CreateDir(d Directory) error
func FileSize(fn Filename) (int64, error)
```

With Printable functions, non-printable runes can be removed from strings and runes

```go
func Printable(a string) string
func IsPrintable(a []string) (bool, error)
func RuneToPrintable(r rune) string
```

With golden file functions, golden files can be created and test cases evaluated. Golden files can be used in unit tests. The expected output is stored in a golden file. The actual output data will be compared with the golden file. The test fails if there is a difference in actual output and golden file.

```go
func GoldenFilePath(name string) (Filename, error)
func CreateGoldenFile(tc *Testcase) error
func EvalGoldenFile(tc *Testcase) error
```

With normalization functions, new lines in byte slices or strings are normalized to the Unix representation of a new line as line feed LF (0x0A). Therefore, Windows new lines CR LF (0x0D 0x0A) are replaced by Unix new lines LF (0x0A). Also, Mac new lines CR (0x0D) are replaced by Unix new lines LF (0x0A).

```go
func NormNewlinesBytes(i []byte) ([]byte, error)
func NormNewlinesStr(i string) string
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/thorsphere/tsfio"
)

func main() {
	f1, _ := os.CreateTemp("", "foo")
	fn1 := tsfio.Filename(f1.Name())

	f2, _ := os.CreateTemp("", "foo")
	fn2 := tsfio.Filename(f2.Name())

	tsfio.RemoveFile(fn1)
	b, _ := tsfio.ExistsFile(fn1)
	fmt.Println(b)

	tsfio.TouchFile(fn1)
	b, _ = tsfio.ExistsFile(fn1)
	fmt.Println(b)

	tsfio.WriteStr(fn1, "foo")
	tsfio.WriteStr(fn1, "foo")
	c, _ := tsfio.ReadFile(fn1)
	fmt.Println(string(c))

	tsfio.WriteStr(fn2, "foo")
	tsfio.AppendFile(&tsfio.Append{FileA: fn1, FileI: fn2})
	c, _ = tsfio.ReadFile(fn1)
	fmt.Println(string(c))

	fs, _ := tsfio.FileSize(fn1)
	fmt.Println(fs)

	tsfio.WriteSingleStr(fn1, "foo")
	c, _ = tsfio.ReadFile(fn1)
	fmt.Println(string(c))

	tsfio.ResetFile(fn1)
	c, _ = tsfio.ReadFile(fn1)
	fmt.Println(string(c))
}
```
[Go Playground](https://go.dev/play/p/FYKm59cQuAV)

Output
```
false
true
foofoo
foofoofoo
9
foo

```

## Links

[Godoc](https://pkg.go.dev/github.com/thorsphere/tsfio)

[Go Report Card](https://goreportcard.com/report/github.com/thorsphere/tsfio)

[Open Source Insights](https://deps.dev/go/github.com%2Fthorsphere%2Ftsfio)
