// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio_test

// Import standard library packages as well as tserr and tsfio
import (
	"fmt"           // fmt
	"os"            // os
	"path/filepath" // filepath
	"testing"       // testing
	"time"          // time

	"github.com/thorsphere/tserr" // tserr
	"github.com/thorsphere/tsfio" // tsfio
)

// TestOpenFile1 tests OpenFile to open an existing temporary
// file and to return *os.File. If OpenFile returns an error or if *os.File is nil,
// the test fails.
func TestOpenFile1(t *testing.T) {
	testOpenFile(t, false)
}

// TestOpenFile2 tests OpenFile to open a new temporary file, not
// existing yet and to return *os.File. If OpenFile returns an error or
// if *os.File is nil, the test fails.
func TestOpenFile2(t *testing.T) {
	testOpenFile(t, true)
}

// TestOpenFile3 tests OpenFile for a provided file name in a not existing directory.
// The test fails if OpenFile does not return an error.
func TestOpenFile3(t *testing.T) {
	// Retrieve a temporary directory
	dn := tmpDir(t)
	// Remove the temporary directory
	rm(t, dn)
	// Retrieve a temporary file
	fn := tmpFile(t)
	// Remove the temporary file
	rm(t, fn)
	// Create a file name with the temporary directory and temporary file name
	fn = tsfio.Filename(filepath.Base(string(fn)))
	n, e := tsfio.Path(dn, fn)
	// The test fails if Path returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Path", Fn: string(dn) + string(fn), Err: e}))
	}
	// Try to open the created file name
	_, e = tsfio.OpenFile(n)
	// The test fails if OpenFile does not return an error
	if e == nil {
		t.Error(tserr.NilFailed("OpenFile"))
	}
}

// testOpenFile is called by Test functions to test OpenFile. If r is true,
// OpenFile opens a non-existing, new temporary file. If r is false,
// OpenFile opens an existing temporary file. If OpenFile returns an error or
// if *os.File is nil, the test fails.
func testOpenFile(t *testing.T, r bool) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create a temporary file with Filename fn
	fn := tmpFile(t)
	// If r is true, remove file fn
	if r {
		rm(t, fn)
	}
	// Open file fn
	f, err := tsfio.OpenFile(fn)
	// If OpenFile returns an error, the test fails
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(fn), Err: err}))
	}
	// If *os.File is nil, the test fails
	if f == nil {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "OpenFile", Actual: "nil", Want: "*os.File"}))
	}
	// Close file fn
	if e := tsfio.CloseFile(f); e != nil {
		// If CloseFile fails, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "CloseFile", Fn: string(fn), Err: e}))
	}
	// Remove file fn
	rm(t, fn)
}

// TestOpenFileEmpty tests OpenFile in case it retrieves the empty string "".
// The test fails if OpenFile returns nil instead of an error.
func TestOpenFileEmpty(t *testing.T) {
	// If OpenFile returns nil, the test fails
	if _, err := tsfio.OpenFile(""); err == nil {
		t.Error(tserr.NilFailed("OpenFile"))
	}
}

// TestCreateDirEmpty tests CreateDor in case it retrieves the empty string "".
// The test fails if CreateDir returns nil instead of an error.
func TestCreateDirEmpty(t *testing.T) {
	// If CreateDir returns nil, the test fails
	if err := tsfio.CreateDir(""); err == nil {
		t.Error(tserr.NilFailed("CreateDir"))
	}
}

// TestCreateDir1 tests CreateDir to create a temporary directory. The test fails,
// if CreateDir returns a Directory that does not exist or if CreateDir returns
// an error.
func TestCreateDir1(t *testing.T) {
	// Create temporary Directory d
	d := tmpDir(t)
	// Remove temporary Directory d
	rm(t, d)
	// Create Directory d with CreateDir
	if err := tsfio.CreateDir(d); err != nil {
		// If CreateDir returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "CreateDir", Fn: string(d), Err: err}))
	}
	// Retrieve FileInfo from Directory d with Stat
	_, e := os.Stat(string(d))
	// If Stat returns that d does not exist, the test fails
	if os.IsNotExist(e) {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "CreateDir", Fn: string(d), Err: e}))
	}
	// In case of any other error of Stat, the test fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "FileInfo (Stat) of", Fn: string(d), Err: e}))
	}
	// Remove Directory d
	rm(t, d)
}

// TestCreateDir2 tests if CreateDir returns nil, in case it retrieves
// a Directory d which already exists. If CreateDir returns an error,
// the test fails.
func TestCreateDir2(t *testing.T) {
	// Create temporary Directory d
	d := tmpDir(t)
	// Call CreateDir with Directory d, which already exists
	if err := tsfio.CreateDir(d); err != nil {
		// If CreateDir returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "CreateDir", Fn: string(d), Err: err}))
	}
	// Remove Directory d
	rm(t, d)
}

// TestCloseFileNil tests if CloseFile returns nil, in case it retrieves nil.
// If CloseFile returns nil, the test fails.
func TestCloseFileNil(t *testing.T) {
	// If CloseFile returns nil, the test fails
	if err := tsfio.CloseFile(nil); err == nil {
		t.Error(tserr.NilFailed("CloseFile"))
	}
}

// TestCloseFile tests CloseFile closing a temporary file. If CloseFile
// returns an error, the test fails.
func TestCloseFile(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Open file fn and get *os.File f
	f, e := tsfio.OpenFile(fn)
	// If OpenFile returns an error, the test fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "OpenFile",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Close file f with CloseFile
	if err := tsfio.CloseFile(f); err != nil {
		// If ClosFile returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "CloseFile",
			Fn:  f.Name(),
			Err: err,
		}))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestCloseFileErr tests if CloseFile returns an error in case
// it is called for a file already closed. The test fails if
// CloseFile returns nil.
func TestCloseFileErr(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Open file fn and get *os.File f
	f, e := tsfio.OpenFile(fn)
	// If OpenFile returns an error, the test fails
	if e != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{
			Op:  "OpenFile",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Close f
	f.Close()
	// Close file f using CloseFile
	if err := tsfio.CloseFile(f); err == nil {
		// If CloseFile returns nil, the test fails
		t.Error(tserr.NilFailed("CloseFile"))
	}
	// Remove file fn
	rm(t, fn)
}

// TestWriteStr tests WriteStr to write two times a string into a new temporary file.
// If the contents of the temporary file does not equal the expected contents,
// the test fails.
func TestWriteStr(t *testing.T) {
	// Define how often WriteStr is called in rep
	rep := 2
	// Set expected test string seq to empty string
	seq := ""
	// Create temporary file fn
	fn := tmpFile(t)
	// Iterate as defined by rep
	for i := 0; i < rep; i++ {
		// If writing the testcase string to fn returns an error, the test fails
		if e := tsfio.WriteStr(fn, testcase); e != nil {
			t.Error(tserr.Op(&tserr.OpArgs{
				Op:  fmt.Sprintf("WriteStr %v to file", testcase),
				Fn:  string(fn),
				Err: e,
			}))
		}
		// Extend expected test string by the testcase string
		seq = seq + testcase
	}
	// Read file fn in b
	b, err := os.ReadFile(string(fn))
	// Remove file fn
	rm(t, fn)
	// If ReadFile returns an error, the test fails
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(fn), Err: err}))
	}
	// If actual b does not match expected seq, the test fails
	if string(b) != seq {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(fn), Actual: string(b), Want: seq}))
	}
}

// TestWriteSingleStr tests WriteSingleStr to write two times a string into a temporary file.
// If the contents of the temporary file does not equal the expected contents, the test fails.
func TestWriteSingleStr(t *testing.T) {
	// Define how often WriteSingleStr is called in rep
	rep := 2
	// Create temporary file fn
	fn := tmpFile(t)
	// Iterate as defined by rep
	for i := 0; i < rep; i++ {
		// If writing the testcase string to fn returns an error, the test fails
		if e := tsfio.WriteSingleStr(fn, testcase); e != nil {
			t.Error(tserr.Op(&tserr.OpArgs{
				Op:  fmt.Sprintf("WriteSingleStr %v to file", testcase),
				Fn:  string(fn),
				Err: e,
			}))
		}
	}
	// Read file fn in b
	b, err := os.ReadFile(string(fn))
	// Remove file fn
	rm(t, fn)
	// If ReadFile returns an error, the test fails
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(fn), Err: err}))
	}
	// If actual b does not match expected seq, the test fails
	if string(b) != testcase {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(fn), Actual: string(b), Want: testcase}))
	}
}

// TestWriteStrEmpty tests WriteStr to return an error for an empty filename.
// If WriteStr returns nil, the test fails.
func TestWriteStrEmpty(t *testing.T) {
	// Write string testcase to empty filename
	if e := tsfio.WriteStr("", testcase); e == nil {
		// If WriteStr returns nil, the test fails
		t.Error(tserr.NilFailed("WriteStr"))
	}
}

// TestWriteSingleStrErr tests WriteSingleStr to return an error for an empty filename.
// If WriteSingleStr returns nil, the test fails.
func TestWriteSingleStrErr(t *testing.T) {
	// Write string testcase to empty filename
	if e := tsfio.WriteSingleStr("", testcase); e == nil {
		// If WriteSingleStr returns nil, the test fails
		t.Error(tserr.NilFailed("WriteSingleStr"))
	}
}

// TestTouchFileEmpty tests TouchFile to return an error for an empty filename.
// If TouchFile returns nil, the test fails.
func TestTouchFileEmpty(t *testing.T) {
	// Touch file with empty filename
	if e := tsfio.TouchFile(""); e == nil {
		// If Touchfile returns nil, the test fails
		t.Error(tserr.NilFailed("TouchFile"))
	}
}

// TestTouchFile1 tests TouchFile to touch an existing temporary
// file. If TouchFile returns an error or if the modification time
// of the temporary file is not later than before, the test fails.
func TestTouchFile1(t *testing.T) {
	testTouchFile(t, false)
}

// TestTouchFile2 tests TouchFile to touch a temporary
// file which was created and removed. If TouchFile returns
// an error or if the modification time of the temporary file
// is not later than before, the test fails.
func TestTouchFile2(t *testing.T) {
	testTouchFile(t, true)
}

// testTouchFile is called by Test functions to test TouchFile. If r is true,
// TouchFile opens a non-existing, new temporary file. If r is false,
// TouchFile opens a temporary file, which was created and removed.
// If TouchFile returns an error or if the modification time of the temporary file
// is not later than before, the test fails.
func testTouchFile(t *testing.T, r bool) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create temporary file fn
	fn := tmpFile(t)
	// Retrieve modification time of fn in t1
	t1 := modTime(t, fn)
	// Remove fn if r is true
	if r {
		rm(t, fn)
	}
	// Pause for at least one second
	time.Sleep(time.Second)
	// Touch file fn
	if e := tsfio.TouchFile(fn); e != nil {
		// If TouchFile returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "TouchFile",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Retrieve modification time of fn in t2
	t2 := modTime(t, fn)
	// Calculate duration d = t2 - t1
	d := t2.Sub(t1)
	// If rounded d is smaller than a second, the test fails
	if d.Round(time.Second) < time.Second {
		t.Error(tserr.Higher(&tserr.HigherArgs{
			Var:        "t2 - t1",
			Actual:     int64(d),
			LowerBound: int64(time.Second),
		}))
	}
	// Remove file fn
	rm(t, fn)
}

// TestReadFile tests ReadFile to read from a temporary file. If the returned
// contents of the temporary file does not match the expected contents,
// the test fails.
func TestReadFile(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Write string testcase to file fn by calling WriteStr
	if e := tsfio.WriteStr(fn, testcase); e != nil {
		// If WriteStr returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  fmt.Sprintf("WriteStr %v to file", testcase),
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Read fn in b
	b, err := tsfio.ReadFile(fn)
	// Remove temporary file fn
	rm(t, fn)
	// If ReadFile returns an error, the test fails
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(fn), Err: err}))
	}
	// If b does not match testcase, the test fails
	if string(b) != testcase {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(fn), Actual: string(b), Want: testcase}))
	}
}

// TestReadFileEmpty tests ReadFile to return an error for an empty filename.
// If ReadFile returns nil, the test fails.
func TestReadFileEmpty(t *testing.T) {
	// Read file with empty Filename
	if _, e := tsfio.ReadFile(""); e == nil {
		// If ReadFile returns nil, the test fails
		t.Error(tserr.NilFailed("ReadFile"))
	}
}

// TestReadFileErr tests ReadFile to return an error when reading a file
// which does not exist. If ReadFile returns nil, the test fails.
func TestReadFileErr(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Write string testcase to fn
	if e := tsfio.WriteStr(fn, testcase); e != nil {
		// If WriteStr returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  fmt.Sprintf("WriteStr %v to file", testcase),
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Remove file fn
	rm(t, fn)
	// Read from file fn
	_, err := tsfio.ReadFile(fn)
	// If ReadFile returns nil, the test fails
	if err == nil {
		t.Error(tserr.NilFailed("ReadFile"))
	}
}

// TestAppendFileNil tests AppendFile to return an error if retrieving a nil pointer.
// If AppendFile returns nil, the test fails.
func TestAppendFileNil(t *testing.T) {
	// Call AppendFile with a nil pointer
	if e := tsfio.AppendFile(nil); e == nil {
		// If AppendFile returns nil, the test fails
		t.Error(tserr.NilFailed("AppendFile"))
	}
}

// TestCopyFileNil tests CopyFile to return an error if retrieving a nil pointer.
func TestCopyFileNil(t *testing.T) {
	// Call CopyFile with a nil pointer
	if e := tsfio.CopyFile(nil); e == nil {
		// If CopyFile returns nil, the test fails
		t.Error(tserr.NilFailed("CopyFile"))
	}
}

// TestAppendFileEmptyA tests AppendFile to return an error if Append argument fileA is
// an empty string. If AppendFile returns nil, the test fails.
func TestAppendFileEmptyA(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Call AppendFile with an empty string as fileA and fn as fileI
	if e := tsfio.AppendFile(&tsfio.Append{FileA: "", FileI: fn}); e == nil {
		// If AppendFile returns nil, the test fails
		t.Error(tserr.NilFailed("AppendFile"))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestCopyFileEmptySrc tests CopyFile to return an error if Copy argument Src is
// an empty string. If CopyFile returns nil, the test fails.
func TestCopyFileEmptySrc(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Call CopyFile with an empty string as source and fn as destination
	if e := tsfio.CopyFile(&tsfio.Copy{Src: "", Dst: fn}); e == nil {
		// If CopyFile returns nil, the test fails
		t.Error(tserr.NilFailed("CopyFile"))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestCopyFileEmptyDst tests CopyFile to return an error if Copy argument Dst is
// an empty string. If CopyFile returns nil, the test fails.
func TestCopyFileEmptyDst(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Call CopyFile with an empty string as destination and fn as source
	if e := tsfio.CopyFile(&tsfio.Copy{Src: fn, Dst: ""}); e == nil {
		// If CopyFile returns nil, the test fails
		t.Error(tserr.NilFailed("CopyFile"))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestAppendFileEmptyI tests AppendFile to return an error if Append argument fileI is
// an empty string. If AppendFile returns nil, the test fails.
func TestAppendFileEmptyI(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Call AppendFile with an empty string as fileI and fn as fileA
	if e := tsfio.AppendFile(&tsfio.Append{FileI: "", FileA: fn}); e == nil {
		// If AppendFile returns nil, the test fails
		t.Error(tserr.NilFailed("AppendFile"))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestAppendFileEmptyIA tests AppendFile to return an error if Append arguments fileI and fileA
// are both an empty string. If AppendFile returns nil, the test fails.
func TestAppendFileEmptyIA(t *testing.T) {
	// Call AppendFile with an empty string as fileI and fileA
	if e := tsfio.AppendFile(&tsfio.Append{FileI: "", FileA: ""}); e == nil {
		// If AppendFile returns nil, the test fails
		t.Error(tserr.NilFailed("AppendFile"))
	}
}

// TestCopyFile tests CopyFile to copy a temporary file to another temporary file.
// If the contents of the destination file does not match the expected contents, the test fails.
func TestCopyFile(t *testing.T) {
	// Create a temporary file as source file
	fnSrc := tmpFile(t)
	// Write the testcase string to the source file
	if e := tsfio.WriteStr(fnSrc, testcase); e != nil {
		// If WriteStr returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  fmt.Sprintf("WriteStr %v to file", testcase),
			Fn:  string(fnSrc),
			Err: e,
		}))
	}

	// Create a temporary file as destination file
	fnDst := tmpFile(t)
	// Remove the destination file to ensure that it does not exist
	rm(t, fnDst)

	// Copy the source file to the destination file using CopyFile
	if e := tsfio.CopyFile(&tsfio.Copy{Src: fnSrc, Dst: fnDst}); e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  fmt.Sprintf("CopyFile %v to file", fnSrc),
			Fn:  string(fnDst),
			Err: e,
		}))
	}

	// Read the contents of the destination file into b using os.ReadFile
	b, err := os.ReadFile(string(fnDst))

	// Remove the source and destination files
	rm(t, fnSrc)
	rm(t, fnDst)

	// If os.ReadFile returns an error, the test fails
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{
			Op:  "ReadFile",
			Fn:  string(fnDst),
			Err: err,
		}))
	}

	// If b does not match the expected string, the test fails
	if string(b) != testcase {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(fnDst), Actual: string(b), Want: testcase}))
	}
}

// TestCopyFileDstExists tests CopyFile to return an error if the destination file already exists.
func TestCopyFileDstExists(t *testing.T) {
	// Create a temporary file as source file
	fnSrc := tmpFile(t)
	// Write the testcase string to the source file
	if e := tsfio.WriteStr(fnSrc, testcase); e != nil {
		// If WriteStr returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  fmt.Sprintf("WriteStr %v to file", testcase),
			Fn:  string(fnSrc),
			Err: e,
		}))
	}

	// Create a temporary file as destination file
	fnDst := tmpFile(t)

	// Copy the source file to the destination file using CopyFile
	if e := tsfio.CopyFile(&tsfio.Copy{Src: fnSrc, Dst: fnDst}); e == nil {
		t.Error(tserr.NilFailed("CopyFile"))
	}

	// Remove the source and destination files
	rm(t, fnSrc)
	rm(t, fnDst)
}

// TestAppendFile tests AppendFile to append a temporary file to another temporary file.
// If the resulting file does not match the expected contents, the test fails.
func TestAppendFile(t *testing.T) {
	// Create two temporary files in fn
	fn := [2]tsfio.Filename{tmpFile(t), tmpFile(t)}
	// Write testcase string to each temporary file in fn
	for _, i := range fn {
		// If WriteStr returns an error, the test fails
		if e := tsfio.WriteStr(i, testcase); e != nil {
			t.Error(tserr.Op(&tserr.OpArgs{
				Op:  fmt.Sprintf("WriteStr %v to file", testcase),
				Fn:  string(i),
				Err: e,
			}))
		}
	}
	// Append fn[1] to fn[0] in fn[0]
	if e := tsfio.AppendFile(&tsfio.Append{FileA: fn[0], FileI: fn[1]}); e != nil {
		// If AppendFile returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  fmt.Sprintf("AppendFile %v to file", fn[1]),
			Fn:  string(fn[0]),
			Err: e,
		}))
	}
	// Read contents of fn[0] in b with os.ReadFile
	b, err := os.ReadFile(string(fn[0]))
	// Remove all files of fn
	for _, i := range fn {
		rm(t, i)
	}
	// If os.ReadFile returns an error, the test fails
	if err != nil {
		t.Fatal(tserr.Op(&tserr.OpArgs{
			Op:  "ReadFile",
			Fn:  string(fn[0]),
			Err: err,
		}))
	}
	// If b does not match the expected string, the test fails
	if string(b) != testcase+testcase {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(fn[0]), Actual: string(b), Want: testcase + testcase}))
	}
}

// TestExistsFileEmpty tests ExistsFile to return an error for an empty string as filename.
// If ExistsFile returns the error to be nil, the test fails.
func TestExistsFileEmpty(t *testing.T) {
	// Call ExistsFile for an empty string as filename
	_, e := tsfio.ExistsFile("")
	// If ExistsFile returns the error to be nil, the test fails
	if e == nil {
		t.Error(tserr.NilFailed("ExistsFile"))
	}
}

// TestExistsDirEmpty tests ExistsDir to return an error for an empty string as directory.
// If ExistsDir returns the error to be nil, the test fails.
func TestExistsDirEmpty(t *testing.T) {
	// Call ExistsFile for an empty string as filename
	_, e := tsfio.ExistsDir("")
	// If ExistsFile returns the error to be nil, the test fails
	if e == nil {
		t.Error(tserr.NilFailed("ExistsFile"))
	}
}

// TestExistsFile1 tests ExistsFile for a non-existing file.
// If ExistsFile returns an error or if ExistsFile returns true,
// the test fails.
func TestExistsFile1(t *testing.T) {
	testExistsFile(t, true)
}

// TestExistsFile2 tests ExistsFile for an existing temporary file.
// If ExistsFile returns an error or if ExistsFile returns false,
// the test fails.
func TestExistsFile2(t *testing.T) {
	testExistsFile(t, false)
}

// TestExistsDir1 tests ExistsDir for a non-existing directory.
// If ExistsDir returns an error or if ExistsDir returns true,
// the test fails.
func TestExistsDir1(t *testing.T) {
	testExistsDir(t, true)
}

// TestExistsDir2 tests ExistsDir for an existing temporary directory.
// If ExistsDir returns an error or if ExistsDir returns false,
// the test fails.
func TestExistsDir2(t *testing.T) {
	testExistsDir(t, false)
}

// testExistsFile is called by Test functions to test ExistsFile. If r is true,
// ExistsFile is tested for a non-existing, removed file. If r is false,
// ExistsFile is tested for an existing temporary file. If ExistsFile returns
// an error or if the actual result of ExistsFile does not match the expected
// result, the test fails.
func testExistsFile(t *testing.T, r bool) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create a temporary file fn
	fn := tmpFile(t)
	// Remove fn, if r is true
	if r {
		rm(t, fn)
	}
	// Call ExistsFile for fn
	b, e := tsfio.ExistsFile(fn)
	// If ExistsFile returns an error, the test fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "ExistsFile",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// If the actual result of ExistsFile does not match the expected result,
	// the test fails
	if b == r {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     "ExistsFile",
			Actual: fmt.Sprintf("%t", b),
			Want:   fmt.Sprintf("%t", !b),
		}))
	}
	// Remove fn, if r is true
	if !r {
		rm(t, fn)
	}
}

// testExistsDir is called by Test functions to test ExistsDir. If r is true,
// ExistsDir is tested for a non-existing, removed file. If r is false,
// ExistsDir is tested for an existing temporary file. If ExistsDir returns
// an error or if the actual result of ExistsDir does not match the expected
// result, the test fails.
func testExistsDir(t *testing.T, r bool) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create a temporary directory dn
	dn := tmpDir(t)
	// Remove dn, if r is true
	if r {
		rm(t, dn)
	}
	// Call ExistsDir for fn
	b, e := tsfio.ExistsDir(dn)
	// If ExistsDir returns an error, the test fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "ExistsFile",
			Fn:  string(dn),
			Err: e,
		}))
	}
	// If the actual result of ExistsDir does not match the expected result,
	// the test fails
	if b == r {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     "ExistsDir",
			Actual: fmt.Sprintf("%t", b),
			Want:   fmt.Sprintf("%t", !b),
		}))
	}
	// Remove fn, if r is true
	if !r {
		rm(t, dn)
	}
}

// TestRemoveFileEmpty tests RemoveFile to return an error for an empty string as filename.
// If RemoveFile returns nil, the test fails.
func TestRemoveFileEmpty(t *testing.T) {
	// Call RemoveFile for an empty string as filename
	if err := tsfio.RemoveFile(""); err == nil {
		// If RemoveFile returns nil, the test fails
		t.Error(tserr.NilFailed("RemoveFile"))
	}
}

// TestRemoveDirEmpty tests RemoveDir to return an error for an empty string as directory.
// If RemoveDir returns nil, the test fails.
func TestRemoveDirEmpty(t *testing.T) {
	// Call RemoveDir for an empty string as directory
	if err := tsfio.RemoveDir(""); err == nil {
		// If RemoveDir returns nil, the test fails
		t.Error(tserr.NilFailed("RemoveDir"))
	}
}

// TestRemoveDir1 tests RemoveDir to remove an existing temporary directory.
// If RemoveDir returns an error or the temporary directory still exists after
// calling RemoveDir, the test fails.
func TestRemoveDir1(t *testing.T) {
	// Create the temporary directory dn
	dn := tmpDir(t)
	// Move a temporary file into dn to ensure that dn is not empty.
	fn := tmpFile(t)
	// Move fn into dn with os.Rename
	os.Rename(string(fn), filepath.Join(string(dn), filepath.Base(string(fn))))
	// Remove dn with RemoveDir
	if e := tsfio.RemoveDir(dn); e != nil {
		// If RemoveDir returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "RemoveDir",
			Fn:  string(dn),
			Err: e,
		}))
	}
	// Check if dn still exists with ExistsDir
	b, err := tsfio.ExistsDir(dn)
	// If ExistsDir returns an error, the test fails
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "ExistsDir",
			Fn:  string(dn),
			Err: err,
		}))
	}
	// If dn still exists, the test fails
	if b {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     "ExistsDir",
			Actual: fmt.Sprintf("%t", b),
			Want:   fmt.Sprintf("%t", !b),
		}))
		// Remove directory dn
		rm(t, dn)
	}
}

// TestRemoveFile1 tests RemoveFile to remove an existing temporary file.
// If RemoveFile returns an error or the temporary file still exists after
// calling RemoveFile, the test fails.
func TestRemoveFile1(t *testing.T) {
	// Create the temporary file fn
	fn := tmpFile(t)
	// Remove fn with RemoveFile
	if e := tsfio.RemoveFile(fn); e != nil {
		// If RemoveFile returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "RemoveFile",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Check if fn still exists with ExistsFile
	b, err := tsfio.ExistsFile(fn)
	// If ExistsFile returns an error, the test fails
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "ExistsFile",
			Fn:  string(fn),
			Err: err,
		}))
	}
	// If fn still exists, the test fails
	if b {
		t.Error(tserr.Return(&tserr.ReturnArgs{
			Op:     "ExistsFile",
			Actual: fmt.Sprintf("%t", b),
			Want:   fmt.Sprintf("%t", !b),
		}))
		// Remove file fn
		rm(t, fn)
	}
}

// TestRemoveDir2 tests RemoveDir for an non-existing directory.
// If RemoveDirectory returns nil, the test fails.
func TestRemoveDir2(t *testing.T) {
	// Create the temporary directory dn
	dn := tmpDir(t)
	// Remove the temporary directory dn
	rm(t, dn)
	// Call RemoveDir for dn
	if e := tsfio.RemoveDir(dn); e == nil {
		// If RemoveDir returns nil, the test fails
		t.Error(tserr.NilFailed("RemoveDir"))
	}
}

// TestRemoveFile2 tests RemoveFile for an non-existing file.
// If RemoveFile returns nil, the test fails.
func TestRemoveFile2(t *testing.T) {
	// Create the temporary file fn
	fn := tmpFile(t)
	// Remove the temporary file fn
	rm(t, fn)
	// Call RemoveFile for fn
	if e := tsfio.RemoveFile(fn); e == nil {
		// If RemoveFile returns nil, the test fails
		t.Error(tserr.NilFailed("RemoveFile"))
	}
}

// TestResetFileEmpty tests ResetFile with an empty string as filename.
// If ResetFiles returns nil, it fails.
func TestResetFileEmpty(t *testing.T) {
	// Call ResetFile with an empty string as filename
	if err := tsfio.ResetFile(""); err == nil {
		// If ResetFile returns nil, the test fails
		t.Error(tserr.NilFailed("ResetFile"))
	}
}

func TestResetFile1(t *testing.T) {
	testResetFile(t, false)
}

func TestResetFile2(t *testing.T) {
	testResetFile(t, true)
}

// testResetFile is called by Test functions to test ResetFile. If r is true,
// ResetFile is tested for a non-existing, removed file. If r is false,
// ResetFile is tested for an existing temporary file. If ResetFile returns
// an error, if FileInfo of the tested file fails, or if the tested file has
// a size greater than zero, the test fails.
func testResetFile(t *testing.T, r bool) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create the temporary file fn
	fn := tmpFile(t)
	// Write single string testcase to fn
	if err := tsfio.WriteSingleStr(fn, testcase); err != nil {
		// If WriteSingleStr returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "WriteSingleStr",
			Fn:  string(fn),
			Err: err,
		}))
	}
	// Remove file fn, if r is true
	if r {
		rm(t, fn)
	}
	// Reset file fn
	if err := tsfio.ResetFile(fn); err != nil {
		// If ResetFile returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "ResetFile",
			Fn:  string(fn),
			Err: err,
		}))
	}
	// Retrieve FileInfo fi from fn with os.Stat
	fi, e := os.Stat(string(fn))
	// If os.Stat returns an error, the test fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "FileInfo (Stat) of",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// If fn has a size greater than zero, the test fails
	if fi.Size() > 0 {
		t.Error(tserr.Equal(&tserr.EqualArgs{
			Var:    fmt.Sprintf("Size of %v", fn),
			Actual: fi.Size(),
			Want:   0,
		}))
	}
	// Remove file fn
	rm(t, fn)
}

// TestFileSizeZero tests FileSize to return 0 for empty files. If FileSize returns
// an error or if FileSize of the tested file is non-zero for an empty temporary file
// the test fails.
func TestFileSizeZero(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Retrieve file size fs for fn
	fs, err := tsfio.FileSize(fn)
	// If FileSize returns an error, the test fails
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "FileSize",
			Fn:  string(fn),
			Err: err,
		}))
	}
	// If file size fs is not zero, the test fails
	if fs != 0 {
		t.Error(tserr.Equal(&tserr.EqualArgs{
			Var:    fmt.Sprintf("FileSize of %v", fn),
			Actual: fs,
			Want:   0,
		}))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestFileSize tests FileSize to return a non-zero value for a non-empty temporary file.
// If FileSize returns an error or if FileSize of the tested file is zero or negative
// the test fails.
func TestFileSize(t *testing.T) {
	// Create temporary file fn
	fn := tmpFile(t)
	// Write testcase to file fn
	if e := tsfio.WriteStr(fn, testcase); e != nil {
		// If WriteStr returns an error, the test fails
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "WriteStr",
			Fn:  string(fn),
			Err: e,
		}))
	}
	// Retrieve file size fs for fn
	fs, err := tsfio.FileSize(fn)
	// If FileSize returns an error, the test fails
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{
			Op:  "FileSize",
			Fn:  string(fn),
			Err: err,
		}))
	}
	// If fs is zero or negative, the test fails
	if fs <= 0 {
		t.Error(tserr.Higher(&tserr.HigherArgs{
			Var:        fmt.Sprintf("FileSize of %v", fn),
			Actual:     fs,
			LowerBound: 1,
		}))
	}
	// Remove temporary file fn
	rm(t, fn)
}

// TestFileSize tests FileSize to return 0 and an error for a temporary directory.
// If FileSize returns nil or if FileSize of the tested directory does not return zero
// the test fails.
func TestFileSizeDir(t *testing.T) {
	// Create a temporary directory d
	d := tmpDir(t)
	// Retrieve file size of d
	fs, err := tsfio.FileSize(tsfio.Filename(d))
	// If FileSize returns nil instead of an error, the test fails
	if err == nil {
		t.Error(tserr.NilFailed(fmt.Sprintf("FileSize of directory %v", d)))
	}
	// If FileSize returns a non-zero size, the test fails
	if fs != 0 {
		t.Error(tserr.Equal(&tserr.EqualArgs{
			Var:    fmt.Sprintf("FileSize of directory %v", d),
			Actual: fs,
			Want:   0,
		}))
	}
	// Remove temporary directory d
	rm(t, d)
}
