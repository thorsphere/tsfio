// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio_test

// Import standard library packages as well as tserr and tsfio
import (
	"os"      // os
	"testing" // testing
	"time"    // time

	"github.com/thorsphere/tsfio"  // tsfio
	"github.com/thorstenrie/tserr" // tserr
)

// The testcases use these tokens
const (
	testprefix    string          = "tsfio_*"         // mostly used as prefix for temporary files or directories
	testcase      string          = "test1234"        // test string
	testcase_unix string          = testcase + "\n"   // test string with a non-printable rune for a new line
	testcase_win  string          = testcase + "\r\n" // test string with a non-printable rune for a new line in Windows
	testcase_mac  string          = testcase + "\r"   // test string with a non-printable rune for a new line in Mac
	testfile      tsfio.Filename  = "test1234"        // test Filename
	testdir       tsfio.Directory = "test/"           // test Directory
	testRP        rune            = 'ú'               // Test printable rune
	testRNp       rune            = '\u001F'          // Test non-printable rune
)

// tmpDir creates a new temporary directory in the default directory for temporary files
// with the prefix testprefix and a random string to the end. In case of an error
// the execution stops.
func tmpDir(t *testing.T) tsfio.Directory {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create the temporary directory
	d, err := os.MkdirTemp("", testprefix)
	// Stop execution in case of an error
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "create temp dir", Fn: d, Err: err}))
	}
	// Return the temporary Directory
	return tsfio.Directory(d)
}

// tmpFile creates a new tmporary file in the default directory for temporary files
// with the prefix testprefix and a random string to the end. In case of an error
// the execution stops.
func tmpFile(t *testing.T) tsfio.Filename {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create the temporary file
	f, err := os.CreateTemp("", testprefix)
	// Stop execution in case of an error
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "create temp file", Fn: f.Name(), Err: err}))
	}
	// Retrieve the name of the temporary file
	fn := tsfio.Filename(f.Name())
	// Close the temporary file
	if e := f.Close(); e != nil {
		// The test fails if Close returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "close temp file", Fn: string(fn), Err: e}))
	}
	// Return the temporary Filename
	return fn
}

// rm removes file named Filename a or empty directory Directory a. In case of an error
// execution stops.
func rm[T tsfio.Fio](t *testing.T, a T) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Remove file or empty directory
	if err := os.Remove(string(a)); err != nil {
		// Stop execution in case of an error
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "Remove", Fn: string(a), Err: err}))
	}
}

// modTime returns the modification time of the file with Filename fn.
// In case of an error, it stops execution.
func modTime(t *testing.T, fn tsfio.Filename) time.Time {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Retrieve the FileInfo structure from fn in fi
	fi, err := os.Stat(string(fn))
	// Stop execution in case of an error
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{
			Op:  "Stat",
			Fn:  string(fn),
			Err: err,
		}))
	}
	// Retrieve modification time from FileInfo fi
	t1 := fi.ModTime()
	// Return the modification time of file with Filename fn
	return t1
}
