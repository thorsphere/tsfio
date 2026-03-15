// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio_test

// Import standard library packages as well as tserr and tsfio
import (
	"testing" // testing

	"github.com/thorsphere/tsfio"  // tsfio
	"github.com/thorstenrie/tserr" // tserr
)

// TestPrintable1 tests Printable, if it returns the identical string, if the provided string a only
// contains printable runes. The test fails, if not.
func TestPrintable1(t *testing.T) {
	// Retrieve the return string from Printable for a test string with only printable runes
	str := tsfio.Printable(testcase)
	// The test fails, if the retrived string does not equal the provided string with only printable runes
	if str != testcase {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "Printable", Actual: str, Want: testcase}))
	}
}

// TestPrintable2 tests Printable to return a string, if the provided string contains a non-printable rune.
// The test fails, if the retrieved string does not equal the provided string without the non-printable rune.
func TestPrintable2(t *testing.T) {
	// Retrieve the return string from Printable for a test string with one non-printable rune
	str := tsfio.Printable(testcase_unix)
	// Return an error if the returned string does not equal the provided string without the non-printable rune
	if str != testcase {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "Printable", Actual: str, Want: testcase}))
	}
}

// TestIsPrintable1 tests IsPrintable to return true, if the provided slice of strings only contains
// printable runes. The test fails, if IsPrintable returns an error or false.
func TestIsPrintable1(t *testing.T) {
	// Retrieve a slice of two strings with only printable runes
	str := []string{testcase, testcase}
	// Retrieve return values from IsPrintable
	b, e := tsfio.IsPrintable(str)
	// The test fails, if IsPrintable returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "IsPrintable", Fn: "slice of strings", Err: e}))
	}
	// The test fails, if IsPrintable returns false
	if !b {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "IsPrintable", Actual: "false", Want: "true"}))
	}
}

// TestIsPrintable2 tests IsPrintable to return false, if the provided slice of strings contains
// a non-printable rune. The test fails, if IsPrintable returns an error or true.
func TestIsPrintable2(t *testing.T) {
	// Retrieve a slice of two strings with one string containing a non-printable rune
	str := []string{testcase_unix, testcase}
	// Retrieve return values of IsPrintable
	b, e := tsfio.IsPrintable(str)
	// The test fails, if IsPrintable returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "IsPrintable", Fn: "slice of strings", Err: e}))
	}
	// The test fails, if IsPrintable returns true
	if b {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "IsPrintable", Actual: "true", Want: "false"}))
	}
}

// TestIsPrintableEmpty1 tests IsPrintable to return false and an error in case of an empty slice of strings. It fails
// if IsPrintable returns nil or true.
func TestIsPrintableEmpty1(t *testing.T) {
	// Retrieve an empty slice of strings
	str := []string{}
	// Test IsPrintable to return false and an error in case of an empty slice of strings. It fails
	// if IsPrintable returns nil or true.
	testIsPrintableEmpty(t, str)
}

// TestIsPrintableEmpty2 tests IsPrintable to return false and an error in case of a slice of strings which is nil. It fails
// if IsPrintable returns nil or true.
func TestIsPrintableEmpty2(t *testing.T) {
	// Retrieve a slice of strings which is nil
	var str []string = nil
	// Test IsPrintable to return false and an error. It fails if IsPrintable returns nil or true.
	testIsPrintableEmpty(t, str)
}

// testIsPrintableEmpty tests IsPrintable to return false and an error in case of an empty slice of strings or a slice of strings
// which is nil. It fails if IsPrintable returns nil or true.
func testIsPrintableEmpty(t *testing.T, str []string) {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// The test fails if length of str is higher than zero
	if len(str) > 0 {
		t.Error(tserr.Equal(&tserr.EqualArgs{Var: "len(str)", Actual: int64(len(str)), Want: 0}))
	}
	// Retrieve return values from IsPrintable for str
	b, e := tsfio.IsPrintable(str)
	// The test fails if IsPrintable returns nil instead of an error
	if e == nil {
		t.Error(tserr.NilFailed("IsPrintable"))
	}
	// The test fails if IsPrintable returns true
	if b {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "IsPrintable", Actual: "true", Want: "false"}))
	}
}

// TestRuneToPrintableEmpty tests RuneToPrintable to return an empty string if the given rune r is empty. The test fails,
// if RuneToPrintable returns a non-empty string.
func TestRuneToPrintableEmpty(t *testing.T) {
	// Retrieve empty rune r
	var r rune
	// Retrieve string from RuneToPrintable for r
	s := tsfio.RuneToPrintable(r)
	// The test fails if the retrieved string s is non-empty
	if s != "" {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: "s", Want: "", Actual: s}))
	}
}

// TestRuneToPrintableP tests RuneToPrintable to return the printable rune as string. The test fails, if
// RuneToPrintable returns a string which does not equal the printable rune as string.
func TestRuneToPrintableP(t *testing.T) {
	// Retrieve printable rune as string s
	s := tsfio.RuneToPrintable(testRP)
	// The test fails, if s does not equal the printable rune as string
	if s != string(testRP) {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: "s", Want: string(testRP), Actual: s}))
	}
}

// TestRuneToPrintableNp tests RuneToPrintable to return an empty string for a non-printable rune. The test fails, if
// RuneToPrintable returns a non-empty string.
func TestRuneToPrintableNp(t *testing.T) {
	// Retrieve string s for the non-printable rune
	s := tsfio.RuneToPrintable(testRNp)
	// The test fails, if s is a non-empty string
	if s != "" {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: "s", Want: "", Actual: s}))
	}
}
