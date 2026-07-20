// Package tsfio provides a file I/O toolkit that extends the standard library with
// utility functions for common file operations and testing workflows.
//
// # Overview
//
// tsfio offers enhanced file handling with built-in safety features:
//   - File operations validate that paths point to regular files or directories only.
//   - Access to system-critical directories and files is blocked on Linux and Windows
//     (see inval_unix.go and inval_win.go for details).
//   - All errors are returned in JSON format using tserr for structured error handling.
//
// # Default File Modes
//
// File and directory operations use sensible defaults:
//   - Files are opened with read-write, append, and create flags (os.O_RDWR | os.O_APPEND | os.O_CREATE).
//   - File permissions default to 0644 (readable by all, writable by owner).
//   - Directory permissions default to 0755 (full access for owner, read/execute for others).
//
// # Features
//
// Printable: Filter non-printable characters from strings and runes.
//
// Golden Files: Manage gold standard test outputs for regression testing. Store expected
// output in golden files and compare against actual results to detect unintended changes.
//
// Normalization: Standardize line endings across platforms. Converts Windows (CR LF) and
// Mac (CR) line endings to Unix line feed format (LF).
//
// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed by the Functional Source License v1.1
// (FSL-1.1-ALv2) that can be found in the LICENSE file.
package tsfio

// Import standard library packages and tserr
import (
	"fmt"           // fmt
	"io/fs"         // fs
	"os"            // os
	"path/filepath" // filepath
	"time"          // time

	"github.com/thorsphere/tserr" // tserr
)

// The constants hold the default flags and file mode
const (
	// Int flags holds the default for files opened read-write, created if it does not exist, data is appended
	flags int = os.O_APPEND | os.O_CREATE | os.O_RDWR
	// FileMode fperm holds the default file mode and permission bits
	fperm fs.FileMode = 0644
	// FileMode dperm holds the default directory mode and permissions bits
	dperm fs.FileMode = 0755
)

// Note: All external file input output functions must contain a CheckFile or CheckDir call at the beginning.
// This can't be tested, because a failed test could break the testing environment.

// OpenFile opens the named file fn with default flags and permission bits. If the file does
// not exist, it is created. If the directory to the file does not exist, the file will
// not be created and OpenFile returns an error.
// If opened successfully, the file is returned and can be used for
// file input output and error is nil.
func OpenFile(fn Filename) (*os.File, error) {
	// Return nil and error in case fn contains a blocked directory or filename
	if e := CheckFile(fn); e != nil {
		return nil, tserr.Check(&tserr.CheckArgs{F: string(fn), Err: e})
	}
	// Get directory of filename, if it contains a path
	if dn := Directory(filepath.Dir(string(fn))); dn != "" {
		// Check if directory exists
		b, err := ExistsDir(dn)
		// Return nil and an error if ExistsDir fails
		if err != nil {
			return nil, tserr.Op(&tserr.OpArgs{Op: "ExistsDir", Fn: string(dn), Err: err})
		}
		// Return nil and an error if the directory does not exist
		if !b {
			return nil, tserr.NotExistent("directory " + string(dn))
		}
	}
	// Open file with default flags and permission bits
	f, err := os.OpenFile(string(fn), flags, fperm)
	// In case of an error, return nil and error
	if err != nil {
		return nil, tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(fn), Err: err})
	}
	// If successful, return file for file input output
	return f, nil
}

// CloseFile closes f and f is unusable for file input output. An error is returned
// if f has already been closed.
func CloseFile(f *os.File) error {
	// Return error in case f is nil
	if f == nil {
		return tserr.NilPtr()
	}
	// Close f
	if e := f.Close(); e != nil {
		// Return error if not successful
		return tserr.Op(&tserr.OpArgs{Op: "Close", Fn: f.Name(), Err: e})
	}
	// If successful, return nil
	return nil
}

// WriteStr writes string s to file with filename fn. It returns an error, if any. If
// fn exists already, the string will be appended to the file. If fn does not exist,
// it will create fn.
func WriteStr(fn Filename, s string) error {
	// Return an error in case fn contains a blocked directory or filename
	if e := CheckFile(fn); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(fn), Err: e})
	}
	// Open file fn with default flags and permission bits. If the file does
	// not exist, it is created.
	f, err := OpenFile(fn)
	// Return error, if OpenFile fails
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(fn), Err: err})
	}
	// Write string s to file fn
	if _, e := f.WriteString(s); e != nil {
		// On error, close file and return error
		f.Close()
		return tserr.Op(&tserr.OpArgs{Op: "write string to", Fn: string(fn), Err: e})
	}
	// Close file
	if e := f.Close(); e != nil {
		// Return error, if Close fails
		return tserr.Op(&tserr.OpArgs{Op: "Close", Fn: string(fn), Err: e})
	}
	// No error occurred, return nil
	return nil
}

// WriteSingleStr writes a single string s to file fn. If fn exists, it is
// truncated to size zero first. If it does not exist, it is created. It
// returns an error, if any.
func WriteSingleStr(fn Filename, s string) error {
	// Return an error in case fn contains a blocked directory or filename
	if e := CheckFile(fn); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(fn), Err: e})
	}
	// Truncate fn to size zero with ResetFile
	if e := ResetFile(fn); e != nil {
		// Return error if ResetFile fails
		return tserr.Op(&tserr.OpArgs{Op: "ResetFile", Fn: string(fn), Err: e})
	}
	// Write string s to fn with WriteStr
	if e := WriteStr(fn, s); e != nil {
		// Return error if WriteStr fails
		return tserr.Op(&tserr.OpArgs{Op: fmt.Sprintf("write string %v to", s), Fn: string(fn), Err: e})
	}
	// No error occurred, return nil
	return nil
}

// TouchFile updates the access and modification times of filename fn to the
// current time. If fn does not exist, it is created as an empty file. It returns
// an error, if any.
func TouchFile(fn Filename) error {
	// Return an error in case fn contains a blocked directory or filename
	if e := CheckFile(fn); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(fn), Err: e})
	}
	// Check if file fn exists
	b, erre := ExistsFile(fn)
	// Return error if ExistsFile fails
	if erre != nil {
		return tserr.Op(&tserr.OpArgs{Op: "ExistsFile", Fn: string(fn), Err: erre})
	}
	// If fn exists, update the access and modification times of fn
	if b {
		// Get current time
		t := time.Now().Local()
		// Update access and modification times of fn. Return error, if any.
		if e := os.Chtimes(string(fn), t, t); e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "Chtimes", Fn: string(fn), Err: e})
		}
	} else {
		// If file does not exist, then create fn with OpenFile.
		f, erro := OpenFile(fn)
		// Return error if OpenFile fails.
		if erro != nil {
			return tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(fn), Err: erro})
		}
		// Close fn
		if e := f.Close(); e != nil {
			// Return error, if Close fails.
			return tserr.Op(&tserr.OpArgs{Op: "Close", Fn: string(fn), Err: e})
		}
	}
	// No error occurred and return nil
	return nil
}

// ReadFile reads f and and returns it contents. It returns an error, if any. If a file
// does not exist, it returns an error. If successful, error will be nil.
func ReadFile(f Filename) ([]byte, error) {
	// Return an error in case f contains a blocked directory or filename
	if e := CheckFile(f); e != nil {
		return nil, tserr.Check(&tserr.CheckArgs{F: string(f), Err: e})
	}
	// Read f and return its contents
	b, err := os.ReadFile(string(f))
	// Return b as nil and the retrieved error, if ReadFile fails
	if err != nil {
		return nil, tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(f), Err: err})
	}
	// No error occurred, return contents and error is nil
	return b, nil
}

// Append holds filename fileA, the file to be extended by contents from file I,
// and filename fileI, which is the input file and will be appended to fileA.
type Append struct {
	FileA Filename // fileA to be extended
	FileI Filename // fileI holds the contents to be appended to fileA
}

// AppendFile appends a file to another file. It receives a pointer to struct
// Append which holds fileA, the file to be extended by contents of fileI and
// fileI, the input file, holding the contents to be appended to fileA. As result,
// fileA holds its original content extended by the contents of fileI, and fileI
// remains as before. If fileA does not exist, it is created as empty file and as
// result will hold the contents of fileI. If fileI does not exist, it returns
// an error. AppendFile returns an error, if any.
func AppendFile(a *Append) error {
	// Return error if pointer a is nil.
	if a == nil {
		return fmt.Errorf("nil pointer")
	}
	// Return an error in case fileA contains a blocked directory or filename
	if e := CheckFile(a.FileA); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(a.FileA), Err: e})
	}
	// Return an error in case fileI contains a blocked directory or filename
	if e := CheckFile(a.FileI); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(a.FileI), Err: e})
	}
	// Read contents of fileI
	out, errr := ReadFile(a.FileI)
	// Return error, if any
	if errr != nil {
		return tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(a.FileI), Err: errr})
	}
	// Open fileA. If it does not exist, then create fileA as empty file.
	f, erro := OpenFile(a.FileA)
	// Return error, if any
	if erro != nil {
		return tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(a.FileA), Err: erro})
	}
	// Write contents of fileI to fileA
	if _, e := f.Write(out); e != nil {
		// If Write filas, close fileA and return error
		f.Close()
		return tserr.Op(&tserr.OpArgs{Op: fmt.Sprintf("append file %v to", a.FileI), Fn: string(a.FileA), Err: e})
	}
	// Close fileA
	if e := f.Close(); e != nil {
		// If Close fails, return error
		return tserr.Op(&tserr.OpArgs{Op: "Close", Fn: string(a.FileA), Err: e})
	}
	// No error occurred, return nil
	return nil
}

// ExistsFile returns true if file fn exists, returns false otherwise. It returns false and an error
// if there is any. If fn is a directory, ExistsFile returns false and an error.
func ExistsFile(fn Filename) (bool, error) {
	// Return an error in case fn is empty
	if fn == "" {
		return false, tserr.Empty("filename")
	}
	// Return exists on fn
	return exists[Filename](fn)
}

// ExistsDir returns true if directory dn exists, returns false otherwise. It returns false and an error
// if there is any. If dn is a file, ExistsDir returns false and an error.
func ExistsDir(dn Directory) (bool, error) {
	// Return an error in case dn is empty
	if dn == "" {
		return false, tserr.Empty("directory name")
	}
	// Return exists on dn
	return exists[Directory](dn)
}

// exists retrieves FileInfo of fn using os.Stat. If os.Stat is successful, it returns true. If os.Stat returns an error reporting fn
// does not exist, it returns false. It returns false and an error for any other error of os.Stat.
func exists[T Fio](fn T) (bool, error) {
	// Retrieve FileInfo of fn
	_, err := os.Stat(string(fn))
	// If Stat is successful return true and error as nil
	if err == nil {
		return true, nil
	}
	// If Stat returns an error reporting fn does not exist, return false and error as nil
	if os.IsNotExist(err) {
		return false, nil
	}
	// For any other error of Stat return false and the error
	return false, tserr.Op(&tserr.OpArgs{Op: "FileInfo (Stat) of", Fn: string(fn), Err: err})
}

// RemoveDir removes directory d and all its contents. It returns an error, if there is any.
// If d does not exist, it returns an error. If d is a file, it also returns an error.
func RemoveDir(d Directory) error {
	// Return an error in case d contains a blocked directory or filename
	if e := CheckDir(d); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(d), Err: e})
	}
	// Check if d exists
	b, err := ExistsDir(d)
	// Return an error if ExistsDir fails
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "check if exists", Fn: string(d), Err: err})
	}
	if b {
		// Remove directory named d, if it exists. Return an error, if RemoveAll fails
		if e := os.RemoveAll(string(d)); e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "RemoveAll", Fn: string(d), Err: e})
		}

	} else {
		// Return an error if d does not exist
		return tserr.NotExistent(string(d))
	}
	// Return nil if no error occurred
	return nil
}

// RemoveFile removes file f. It returns an error, if there is any. If f is a directory
// it returns an error. If f does not exist, it also returns an error.
func RemoveFile(f Filename) error {
	// Return an error in case f contains a blocked directory or filename
	if e := CheckFile(f); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(f), Err: e})
	}
	// Check if f exists
	b, err := ExistsFile(f)
	if err != nil {
		// Return an error if ExistsFile fails
		return tserr.Op(&tserr.OpArgs{Op: "check if exists", Fn: string(f), Err: err})
	}
	if b {
		// Remove f, if it exists
		if e := os.Remove(string(f)); e != nil {
			// Return an error if Remove fails
			return tserr.Op(&tserr.OpArgs{Op: "Remove", Fn: string(f), Err: e})
		}
	} else {
		// Return an error if f does not exist
		return tserr.NotExistent(string(f))
	}
	// No error occurred, return nil
	return nil
}

// ResetFile truncates fn to size zero. If fn does not exist, it is created as empty file.
// It returns an error, if there is any.
func ResetFile(fn Filename) error {
	// Return an error in case fn contains a blocked directory or filename
	if e := CheckFile(fn); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(fn), Err: e})
	}
	// Check if fn exists
	b, err := ExistsFile(fn)
	if err != nil {
		// Return error, if ExistsFile fails
		return tserr.Op(&tserr.OpArgs{Op: "check if exists", Fn: string(fn), Err: err})
	}
	if !b {
		// If fn does not exist, it is created as an empty file with TouchFile
		if e := TouchFile(fn); e != nil {
			// Return error, if TouchFile fails
			return tserr.Op(&tserr.OpArgs{Op: "TouchFile", Fn: string(fn), Err: e})
		}
	}
	// Truncate file fn to size zero
	err = os.Truncate(string(fn), 0)
	if err != nil {
		// Return error if Truncate fails
		return tserr.Op(&tserr.OpArgs{Op: "Truncate", Fn: string(fn), Err: err})
	}
	// No error occurred, return nil
	return nil
}

// CreateDir creates a directory named d with any necessary parents. If d already exists as
// directory, it does nothing and returns nil. It returns an error, if there is any.
func CreateDir(d Directory) error {
	// Return an error in case d contains a blocked directory or filename
	if e := CheckDir(d); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(d), Err: e})
	}
	// Create directory named d with any necessary parents
	err := os.MkdirAll(string(d), dperm)
	if err != nil {
		// Return an error, if MKdirAll fails
		return tserr.Op(&tserr.OpArgs{Op: "make directory", Fn: string(d), Err: err})
	}
	// No error occurred, return nil
	return nil
}

// FileSize returns the length in bytes for the regular file fn. If fn is a blocked filename,
// a directory or if FileInfo for fn cannot be retrieved, it returns 0 and an error.
func FileSize(fn Filename) (int64, error) {
	// Return an error in case fn contains a blocked directory or filename
	if err := CheckFile(fn); err != nil {
		return 0, tserr.Check(&tserr.CheckArgs{F: string(fn), Err: err})
	}
	// Retrieve FileInfo of fn
	fi, e := os.Stat(string(fn))
	if e != nil {
		// For any error of Stat return 0 and the error
		return 0, tserr.Op(&tserr.OpArgs{Op: "FileInfo (Stat) of", Fn: string(fn), Err: e})
	}
	return fi.Size(), nil
}

// Copy holds the source file Src and the destination file Dst for the copy operation of CopyFile.
type Copy struct {
	Src Filename // Src holds the source file for the copy operation of CopyFile
	Dst Filename // Dst holds the destination file for the copy operation of CopyFile
}

// CopyFile copies the contents of file src to dst. If dst does not exist,
// it is created as an empty file and as result will hold the contents of src.
// It returns an error if dst already exists, if src does not exist,
// if src or dst contains a blocked directory or filename,
// if src is a directory, if dst is a directory or
// if there is any error during the copy process.
func CopyFile(a *Copy) error {
	// Return error if pointer a is nil.
	if a == nil {
		return tserr.NilFailed("CopyFile")
	}
	// Return an error in case src contains a blocked directory or filename
	if e := CheckFile(a.Src); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(a.Src), Err: e})
	}
	// Return an error in case dst contains a blocked directory or filename
	if e := CheckFile(a.Dst); e != nil {
		return tserr.Check(&tserr.CheckArgs{F: string(a.Dst), Err: e})
	}

	// Open source file
	srcFile, err := os.Open(string(a.Src))
	// Return an error if Open fails
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "Open", Fn: string(a.Src), Err: err})
	}

	// Check if destination file already exists
	b, err := ExistsFile(a.Dst)
	// Return an error if ExistsFile fails
	if err != nil {
		// Attempt to close source file before returning error
		srcFile.Close()
		// Return error if ExistsFile fails
		return tserr.Op(&tserr.OpArgs{Op: "ExistsFile", Fn: string(a.Dst), Err: err})
	}

	// Return an error if destination file already exists
	if b {
		// Attempt to close source file before returning error
		srcFile.Close()
		// Return error if destination file already exists
		return tserr.AlreadyExistent(string(a.Dst))
	}

	// Create destination file
	dstFile, err := OpenFile(a.Dst)
	if err != nil {
		// Attempt to close source file before returning error
		srcFile.Close()
		// Return error if OpenFile fails
		return tserr.Op(&tserr.OpArgs{Op: "Create", Fn: string(a.Dst), Err: err})
	}

	// Copy contents from source to destination
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		// Attempt to close files before returning error
		srcFile.Close()
		dstFile.Close()
		// Return error if ReadFrom fails
		return tserr.Op(&tserr.OpArgs{Op: "ReadFrom", Fn: fmt.Sprintf("%v to %v", a.Src, a.Dst), Err: err})
	}

	// Close source file
	if err := srcFile.Close(); err != nil {
		// Attempt to close destination file before returning error
		dstFile.Close()
		// Return error if Close fails
		return tserr.Op(&tserr.OpArgs{Op: "Close", Fn: string(a.Src), Err: err})
	}

	// Close destination file
	if err := dstFile.Close(); err != nil {
		// Return error if Close fails
		return tserr.Op(&tserr.OpArgs{Op: "Close", Fn: string(a.Dst), Err: err})
	}

	// No error occurred, return nil
	return nil
}
