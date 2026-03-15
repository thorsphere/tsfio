// Package tsfio provides a simple API for file input output based on the standard library.
//
// The tsfio package is a supplement to the standard library and supplies additional
// functions for file input output operations, e.g., appending one file to another file.
// Also, file input output operations on Linux and Windows system directories or
// files are blocked (see inval_unix.go and inval_win.go) and an error is returned.
// All operations expect a directory or a regular file, return an error otherwise.
// Default flags and file mode is used when opening files, creating files or directories
// and when writing to files (with exceptions documented in the function descriptions)
//
//   - Files are opened read-write (os.O_RDWR).
//   - Data is appended when writing to file (os.O_APPEND).
//   - A file is created if it does not exist (os.O_CREATE).
//   - File mode and permission bits are 0644.
//   - Directory mode and permissions bits are 0755.
//
// If an API call is not successful, a tserr error in JSON format is returned.
//
// With Printable functions, non-printable runes can be removed from strings and runes.
// With golden file functions, golden files can be created and test cases evaluated.
// Golden files can be used in unit tests. The expected output is stored in a golden file.
// The actual output data will be compared with the golden file. The test fails if there
// is a difference in actual output and golden file. With normalization functions, new lines
// in byte slices or strings are normalized to the Unix representation of a new line as
// line feed LF (0x0A). Therefore, Windows new lines CR LF (0x0D 0x0A) are replaced by Unix
// new lines LF (0x0A). Also, Mac new lines CR (0x0D) are replaced by Unix new lines LF (0x0A).
//
// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio

// Import standard library packages and tserr
import (
	"fmt"           // fmt
	"io/fs"         // fs
	"os"            // os
	"path/filepath" // filepath
	"time"          // time

	"github.com/thorstenrie/tserr" // tserr
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
		e := os.Remove(string(f))
		if e != nil {
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
