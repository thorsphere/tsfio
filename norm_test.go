// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio_test

// Import standard library packages as well as tserr and tsfio
import (
	"bytes"   // bytes
	"testing" // testing

	"github.com/thorsphere/tserr" // tserr
	"github.com/thorsphere/tsfio" // tsfio
)

// testNormByte tests NormNewlinesBytes to return normalized new lines for input string i with name n. New lines are expected to be normalized
// to the Unix representation of a new line as line feed LF (0x0A). The function is intended to be used with variables testcase_{unix, win, mac}.
// If the normalized string does not equal ref, the test fails.
func testNormByte(i []byte, n, ref string, t *testing.T) {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// Convert ref to a byte slice
	w := []byte(ref)
	// Retrieve normalized byte slice for i
	test, err := tsfio.NormNewlinesBytes(i)
	// The test fails if NormNewlinesBytes returns an error
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "NormNewlinesBytes", Fn: n, Err: err}))
	}
	// The test fails if the retrieved normalized byte slice does not equal ref
	if !bytes.Equal(test, w) {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: n, Actual: string(test), Want: ref}))
	}
}

// testNormStr tests NormNewlinesStr to return normalized new lines for input string i with name n. New lines are expected to be normalized
// to the Unix representation of a new line as line feed LF (0x0A). The function is intended to be used with variables testcase_{unix, win, mac}.
// If the normalized string does not equal ref, the test fails.
func testNormStr(i, n, ref string, t *testing.T) {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// Retrieve normalized string from i
	test := tsfio.NormNewlinesStr(i)
	// The test fails if the retrieved normalized string does not equal ref
	if test != ref {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: n, Actual: test, Want: ref}))
	}
}

// TestNormByte tests NormNewlinesBytes to return normalized new lines for input byte slices with Unix, Windows and Mac line endings. New lines are expected to be normalized
// to the Unix representation of a new line as line feed LF (0x0A). The function is intended to be used with variables testcase_{unix, win, mac}.
// If the normalized byte slice does not equal testcase_unix, the test fails.
func TestNormByte(t *testing.T) {
	// Testcase with Unix line ending
	testNormByte([]byte(testcase_unix), "testcase_unix", testcase_unix, t)
	// Testcase with Windows line ending
	testNormByte([]byte(testcase_win), "testcase_win", testcase_unix, t)
	// Testcase with Mac line ending
	testNormByte([]byte(testcase_mac), "testcase_mac", testcase_unix, t)
}

// TestNormStr tests NormNewlinesStr to return normalized new lines for input strings with Unix, Windows and Mac line endings. New lines are expected to be normalized
// to the Unix representation of a new line as line feed LF (0x0A). The function is intended to be used with variables testcase_{unix, win, mac}.
// If the normalized string does not equal testcase_unix, the test fails.
func TestNormStr(t *testing.T) {
	// Testcase with Unix line ending
	testNormStr(testcase_unix, "testcase_unix", testcase_unix, t)
	// Testcase with Windows line ending
	testNormStr(testcase_win, "testcase_win", testcase_unix, t)
	// Testcase with Mac line ending
	testNormStr(testcase_mac, "testcase_mac", testcase_unix, t)
}

// TestNormByteNil tests NormNewlinesBytes to return an error in case the input byte slice is nil.
// The test fails, if NormNewlinesBytes returns nil instead of an error.
func TestNormByteNil(t *testing.T) {
	// Retrieve o and err from NormNewlinesBytes in case i is nil
	o, err := tsfio.NormNewlinesBytes(nil)
	// The test fails if NormNewlinesBytes returns nil instead of an error or an output o which is not nil.
	if (err == nil) || (o != nil) {
		t.Error(tserr.NilFailed("NormNewlinesBytes"))
	}
}
