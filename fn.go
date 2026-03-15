// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio

// Import standard library packages and tserr
import (
	"fmt"           // fmt
	"os"            // os
	"path/filepath" // path/filepath
	"strings"       // strings

	"github.com/thorstenrie/tserr" // tserr
)

// Interface Fio is constrained to type Filename and Directory
type Fio interface {
	Filename | Directory
}

// A Filename is the name of a regular file and may contain its path
type Filename string

// A Directory is the name of a directory and may contain its path
type Directory string

// CheckFile performs checks on file f and returns an error if
//   - f is an empty string
//   - f contains a blocked directory or filename
//   - f is an existing directory, not a file
//   - os.Stat returns an error when retrieving FileInfo
//
// Otherwise it returns nil.
func CheckFile(f Filename) error {
	return checkWrapper(f, false)
}

// CheckDir performs checks on directory d and returns an error if
//   - d is an empty string
//   - d contains a blocked directory or filename
//   - d is an existing file, not a directory
//   - os.Stat returns an error when retrieving FileInfo
//
// Otherwise it returns nil.
func CheckDir(d Directory) error {
	return checkWrapper(d, true)
}

// checkWrapper performs checks on a file or directory using fio as type parameter.
// It returns an error, if any check fails. Otherwise it returns nil.
func checkWrapper[T Fio](f T, dir bool) error {
	// Return an error if f is an empty string
	if f == "" {
		return tserr.Empty(string(f))
	}
	// Return an error if f contains a blocked directory or filename
	if err := checkInval(f); err != nil {
		return err
	}
	// Retrieve FileInfo of f
	i, err := os.Stat(string(f))
	// Set w to file or directory
	w := "regular file"
	if dir {
		w = "directory"
	}
	// If Stat returns no error, then check if expected type matches type of f
	if err == nil {
		// Return nil, if expected type matches
		if i.IsDir() && dir {
			return nil
		}
		if i.Mode().IsRegular() && !dir {
			return nil
		}
		// Return an error otherwise
		return tserr.TypeNotMatching(&tserr.TypeNotMatchingArgs{Actual: string(f), Want: w})
	}
	// If Stat returns an error reporting f does not exist, return nil
	if os.IsNotExist(err) {
		return nil
	}
	// If Stat returns any other error, return the error
	return tserr.Check(&tserr.CheckArgs{F: string(f), Err: err})
}

// checkInval if f contains blocked directories or equals a blocked filename.
// In case of a match with a blocked directory or filename it returns an error,
// otherwise nil.
func checkInval[T Fio](f T) error {
	// Retrieve the shortest path name of f
	fc := filepath.Clean(string(f))
	// Iterate i over blocked filenames
	for _, i := range InvalFile() {
		// Retrieve the shortest path name of i
		ic := filepath.Clean(string(i))
		// If the blocked filename and f match, then return an error
		if ic == fc {
			return tserr.Forbidden(string(f))
		}
	}
	// Iterate i over blocked directories
	for _, i := range InvalDir() {
		// Retrieve the shortest path name of i
		ic := filepath.Clean(string(i))
		// If f matches the blocked directory or one of its parents, then return an error
		if strings.HasPrefix(fc, ic) {
			return tserr.Forbidden(string(i))
		}
	}
	// No error occurred, return nil
	return nil
}

// Sprintf formats according to the format specifier and returns the resulting Filename or Directory
func Sprintf[T Fio](f string, a ...any) T {
	return T(fmt.Sprintf(f, a...))
}

// Path joins directory name d and a filename f into a single path p. It returns an empty string and an error if checks
// on d, f and p fail. Path joins the path elements by using the Join function from the Go standard library package path/filepath.
func Path(d Directory, f Filename) (Filename, error) {
	// Return an error in case d contains a blocked directory or filename
	if e := CheckDir(d); e != nil {
		return "", tserr.Check(&tserr.CheckArgs{F: string(d), Err: e})
	}
	// Return an error in case f contains a blocked directory or filename
	if e := CheckFile(f); e != nil {
		return "", tserr.Check(&tserr.CheckArgs{F: string(f), Err: e})
	}
	p := Filename(filepath.Join(string(d), string(f)))
	// Return an error in case p contains a blocked directory or filename
	if e := CheckFile(p); e != nil {
		return "", tserr.Check(&tserr.CheckArgs{F: string(p), Err: e})
	}
	return p, nil
}
