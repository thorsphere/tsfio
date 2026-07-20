// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed by the Functional Source License v1.1
// (FSL-1.1-ALv2) that can be found in the LICENSE file.
package tsfio_test

// Import standard library packages as well as tserr and tsfio
import (
	"testing" // testing

	"github.com/thorsphere/tserr" // tserr
	"github.com/thorsphere/tsfio" // tsfio
)

// TestGoldenFile1 tests the creation and evaluation of golden files. It fails in case of an error.
func TestGoldenFile1(t *testing.T) {
	testGoldenFile(t, true)
}

// TestGoldenFile2 tests the creation and evaluation of golden files. It tests a golden file and a test case not being
// equal. It fails if the evaluation of the golden file does not return an error.
func TestGoldenFile2(t *testing.T) {
	testGoldenFile(t, false)
}

// testGoldenFile tests the creation and evaluation of golden files. If r is true, the testcase and golden file are equal and the evaluation is expected to be
// successful. If r is false, the testcase and golden file are not equal and the evaluation is expected to fail.
func testGoldenFile(t *testing.T, r bool) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create the testcase
	tc := &tsfio.Testcase{Name: testcase, Data: testcase}
	// Create the golden file
	e := tsfio.CreateGoldenFile(tc)
	// The test fails if CreateGoldenFile returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "CreateGoldenFile", Fn: tc.Name, Err: e}))
	}
	// Change testcase, if r is false
	if !r {
		tc.Data += testcase
	}
	// Evaluate the golden file
	e = tsfio.EvalGoldenFile(tc)
	// The test fails if EvalGoldenFile returns an error and an error is not expected
	if r && (e != nil) {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "EvalGoldenFile", Fn: tc.Name, Err: e}))
	} else if !r && (e == nil) {
		// The test fails if EvalGoldenFile returns nil and an error is expected
		t.Error(tserr.NilFailed("EvalGoldenFile"))
	}
	// Retrieve the golden file path of the testcase
	fn, e := tsfio.GoldenFilePath(tc.Name)
	// The test fails if GoldenFilePath returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "GoldenFilePath", Fn: tc.Name, Err: e}))
	}
	// Retrieve whether the golden file exists
	b, e := tsfio.ExistsFile(fn)
	// The test fails if ExistsFile returns an error
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "ExistsFile", Fn: string(fn), Err: e}))
	}
	// The test fails if the golden file does not exist
	if !b {
		t.Error(tserr.NotExistent(string(fn)))
	}
	// Remove the golden file
	if e := tsfio.RemoveFile(fn); e != nil {
		// The test fails if RemoveFile returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "RemoveFile", Fn: string(fn), Err: e}))
	}
}

// TestCreateGoldenFileNil tests CreateGoldenFile to return an error if the testcase is nil.
// The test fails if CreateGoldenFile returns nil instead of an error.
func TestCreateGoldenFileNil(t *testing.T) {
	if e := tsfio.CreateGoldenFile(nil); e == nil {
		t.Error(tserr.NilFailed("CreateGoldenFile"))
	}
}

// TestEvalGoldenFileNil tests EvalGoldenFile to return an error if the testcase is nil.
// The test fails if EvalGoldenFile returns nil instead of an error.
func TestEvalGoldenFileNil(t *testing.T) {
	if e := tsfio.EvalGoldenFile(nil); e == nil {
		t.Error(tserr.NilFailed("EvalGoldenFile"))
	}
}
