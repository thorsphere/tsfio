// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio_test

// Import standard library packages as well as tserr and tsfio
import (
	"fmt" // fmt
	"os"  // os
	"path/filepath"
	"testing" // testing

	"github.com/thorsphere/tsfio"  // tsfio
	"github.com/thorstenrie/tserr" // tserr
)

// TestBlockedDir tests if CheckFile returns an error for all
// blocked directories in invalDir. If it returns nil for
// a blocked directory, the test fails.
func TestBlockedDir(t *testing.T) {
	// Iterate test over all directories in invalDir
	for _, d := range tsfio.InvalDir() {
		// Create test Filename p containing the blocked directory
		p := tsfio.Filename(d) + tsfio.Filename(os.PathSeparator) + testfile
		// If CheckFile returns nil, then the test fails
		if tsfio.CheckFile(p) == nil {
			t.Error(tserr.NilFailed(fmt.Sprintf("CheckFile of %v", d)))
		}
	}
}

// TestBlockedFile tests if CheckFile returns an error for all
// blocked files in invalFile. If it returns nil for a blocked
// file, the test fails.
func TestBlockedFile(t *testing.T) {
	// Iterate test over all files in invalFile
	for _, d := range tsfio.InvalFile() {
		// If CheckFile returns nil, then the test fails
		if tsfio.CheckFile(d) == nil {
			t.Error(tserr.NilFailed(fmt.Sprintf("CheckFile of %v", d)))
		}
	}
}

// TestEmptyDir tests if CheckDir returns an error for an empty string
// as Directory. If it returns nil for the empty string as Directory,
// the test fails.
func TestEmptyDir(t *testing.T) {
	// If CheckDir of the empty string as Directory returns nil, then the test fails
	if err := tsfio.CheckDir(""); err == nil {
		t.Error(tserr.NilFailed("CheckDir"))
	}
}

// TestEmptyFile tests if CheckFile returns an error for an empty string
// as Filename. If it returns nil for the empty string as Filename,
// the test fails.
func TestEmptyFile(t *testing.T) {
	// If CheckFile of the empty string as Filename returns nil, then the test fails
	if err := tsfio.CheckFile(""); err == nil {
		t.Error(tserr.NilFailed("CheckFile"))
	}
}

// TestDir1 tests if CheckDir returns nil for a newly created temporary directory.
// If it returns an error, the test fails.
func TestDir1(t *testing.T) {
	// Create a temporary directory with name d
	d := tmpDir(t)
	// Test fails, if Checkdir returns an error for d
	if err := tsfio.CheckDir(d); err != nil {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     fmt.Sprintf("CheckDir of %v", d),
			Actual: fmt.Sprint(err),
			Want:   "nil",
		}))
	}
	// Remove the temporary directory d
	rm(t, d)
}

// TestDir2 tests if CheckDir returns nil for a newly created temporary directory
// which is removed before the check. If it returns an error, the test fails.
func TestDir2(t *testing.T) {
	// Create a temporary directory with name d
	d := tmpDir(t)
	// Remove the temporary directory d before the check
	rm(t, d)
	// Test fails, if Checkdir returns an error for d
	if err := tsfio.CheckDir(d); err != nil {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     fmt.Sprintf("CheckDir of %v", d),
			Actual: fmt.Sprint(err),
			Want:   "nil",
		}))
	}
}

// TestFile1 tests if CheckFile returns nil for a newly created temporary file.
// If it returns an error, the test fails.
func TestFile1(t *testing.T) {
	// Create a temporary file with name f
	f := tmpFile(t)
	// Test fails, if CheckFile returns an error for f
	if err := tsfio.CheckFile(f); err != nil {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     fmt.Sprintf("CheckFile of %v", f),
			Actual: fmt.Sprint(err),
			Want:   "nil",
		}))
	}
	// Remove the temporary file f
	rm(t, f)
}

// TestFile2 tests if CheckFile returns nil for a newly created temporary file
// which is removed before the check. If it returns an error, the test fails.
func TestFile2(t *testing.T) {
	// Create a temporary file with name f
	f := tmpFile(t)
	// Remove the temporary file f before the check
	rm(t, f)
	// Test fails, if CheckFile returns an error for f
	if err := tsfio.CheckFile(f); err != nil {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     fmt.Sprintf("CheckFile of %v", f),
			Actual: fmt.Sprint(err),
			Want:   "nil",
		}))
	}
}

// TestDir3 tests if CheckDir returns an error when it receives
// a temporary file. If it returns nil, the test fails.
func TestDir3(t *testing.T) {
	// Create a temporary file with name f
	f := tmpFile(t)
	// If CheckDir returns nil, the test fails
	if tsfio.CheckDir(tsfio.Directory(f)) == nil {
		t.Error(tserr.NilFailed(fmt.Sprintf("CheckDir of file %v", f)))
	}
	// Remove the temporary file f
	rm(t, f)
}

// TestFile3 tests if CheckFile returns an error when it receives a
// temporary directory. If it returns nil, the test fails.
func TestFile3(t *testing.T) {
	// Create a temporary directory with name d
	d := tmpDir(t)
	// If CheckFile returns nil, the test fails
	if tsfio.CheckFile(tsfio.Filename(d)) == nil {
		t.Error(tserr.NilFailed(fmt.Sprintf("CheckFile of directory %v", d)))
	}
	// Remove the temporary directory d
	rm(t, d)
}

// TestSprintf tests if Sprintf formats according to the format specifier and
// returns the resulting Filename or Directory. The test fails if not.
func TestSprintf(t *testing.T) {
	// Wanted string swant
	swant := os.TempDir() + string(os.PathSeparator) + testcase
	// Use Sprintf for a Directory d
	d := tsfio.Sprintf[tsfio.Directory]("%v%v%v", os.TempDir(), string(os.PathSeparator), testcase)
	// Use Sprintf for a Filename f
	f := tsfio.Sprintf[tsfio.Filename]("%v%v%v", os.TempDir(), string(os.PathSeparator), testcase)
	// If d does not equal swant, the test fails
	if tsfio.Directory(swant) != d {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{
			Var:    "d",
			Want:   swant,
			Actual: string(d),
		}))
	}
	// If f does not equal swant, the test fails
	if tsfio.Filename(swant) != f {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{
			Var:    "f",
			Want:   swant,
			Actual: string(f),
		}))
	}
}

// TestEmptyJoin tests Path to return an error if d and f are empty. The test fails if the returned path is not empty or the error is nil.
func TestEmptyJoin(t *testing.T) {
	// Retrieve p and e from Path
	p, e := tsfio.Path("", "")
	// The test fails if the returned path is not empty or the error is nil
	if p != "" || e == nil {
		t.Error(tserr.NilFailed("Path"))
	}
}

// TestBlockedDirJoin tests Path to return an empty path and an error in case d equals a blocked directory. The test
// fails if Path does not return an empty path or the error is nil.
func TestBlockedDirJoin(t *testing.T) {
	// The test fails if InvalDir is empty
	if len(tsfio.InvalDir()) == 0 {
		t.Fatal(tserr.NilPtr())
	}
	// Retrieve first blocked directory in d
	d := tsfio.InvalDir()[0]
	// Retrieve test filename in f
	f := testfile
	// Retrieve p from path with blocked directory d and test filename f
	p, e := tsfio.Path(d, f)
	// The test fails if Path does not return an empty path or the is nil
	if p != "" || e == nil {
		t.Error(tserr.NilFailed("Path"))
	}
}

// TestBlockedFileJoin tests Path to return an empty path and an error in case f equals a blocked filename. The test
// fails if Path does not return an empty path or the error is nil.
func TestBlockedFileJoin(t *testing.T) {
	// The test fails if InvalFile is empty
	if len(tsfio.InvalFile()) == 0 {
		t.Fatal(tserr.NilPtr())
	}
	// Retrieve the test directory name in d
	d := testdir
	// Retrieve first blocked filename in f
	f := tsfio.InvalFile()[0]
	// Retrieve p from Path with test directory d and blocked filename f
	p, e := tsfio.Path(d, f)
	// The test fails if Path does not return an empty path or the is nil
	if p != "" || e == nil {
		t.Error(tserr.NilFailed("Path"))
	}
}

// TestJoin tests Path to return the joined directory d and filename f. The test fails if Path returns
// an error or if Path does not match the expected result.
func TestJoin(t *testing.T) {
	// Retrieve test directory name in d
	d := testdir
	// Retrieve test filename in f
	f := testfile
	// Retrieve p from Path
	p, e := tsfio.Path(d, f)
	// The test fails if Path returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Path", Fn: fmt.Sprintf("%v,%v", d, f), Err: e}))
	}
	// Retrieve expected result from filepath.Join
	r := tsfio.Filename(filepath.Join(string(d), string(f)))
	// The test fails if the result of Path in p does not equal the expected result
	if p != r {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: "p", Actual: string(p), Want: string(r)}))
	}
}

// TestInvalDir tests InvalDir to return copies of the array of blocked directories. It retrieves two copies. One copy of the array
// is changed. The test fails if both copies are equal, if only one copy is changed.
func TestInvalDir(t *testing.T) {
	// Retrieve a copy of blocked directories from InvalDir
	d1 := tsfio.InvalDir()
	// Retrieve a copy of blocked directories from InvalDir
	d2 := tsfio.InvalDir()
	// Change the first copy
	d1[0] = testdir
	// The test fails if both copies are equal
	if d1 == d2 {
		t.Error(tserr.NotEqual(&tserr.NotEqualArgs{X: "d1", Y: "d2"}))
	}
}

// TestInvalFile tests InvalFIle to return copies of the array of blocked filenames. It retrieves two copies. One copy of the array
// is changed. The test fails if both copies are equal, if only one copy is changed.
func TestInvalFile(t *testing.T) {
	// Retrieve a copy of blocked directories from InvalFile
	f1 := tsfio.InvalFile()
	// Retrieve a copy of blocked directories from InvalFile
	f2 := tsfio.InvalFile()
	// Change the first copy
	f1[0] = testfile
	// The test fails if both copies are equal
	if f1 == f2 {
		t.Error(tserr.NotEqual(&tserr.NotEqualArgs{X: "f1", Y: "f2"}))
	}
}
