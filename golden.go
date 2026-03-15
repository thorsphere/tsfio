// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tsfio

// Import package tserr
import "github.com/thorstenrie/tserr" // tserr

// Default directory and file type of golden files
const (
	goldenDir      Directory = "testdata/" // Default golden files directory
	goldenFileType string    = ".golden"   // Default golden file type
)

// A Testcase contains the name of a testcase and the corresponding data of the testcase.
// The data can be reference data or test data.
type Testcase struct {
	Name string // Name of the testcase
	Data string // Data of the testcase
}

// GoldenFilePath returns the path of the test data golden file for the provided name of a testcase.
// The golden files are stored in the default golden files directory testdata/ and have the default
// golden file type .golden.
func GoldenFilePath(name string) (Filename, error) {
	return Path(goldenDir, Filename(name+goldenFileType))
}

// CreateGoldenFile creates a golden file provided by the testcase name. The data in the testcase is written to
// the golden file. The golden file is stored in the default golden files directory testdata/ and has the default
// golden file type .golden.
func CreateGoldenFile(tc *Testcase) error {
	// Return an error if tc is nil
	if tc == nil {
		return tserr.NilPtr()
	}
	if e := CreateDir(goldenDir); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "CreateDir", Fn: string(goldenDir), Err: e})
	}
	// Retrieve golden file path for testcase name
	fn, e := GoldenFilePath(tc.Name)
	// Return an error if goldenPath fails
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "goldenPath", Fn: tc.Name, Err: e})
	}
	// Write the data from the testcase to the golden file
	e = WriteSingleStr(fn, tc.Data)
	// Return an error if WriteSingleStr fails
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "WriteSingleStr", Fn: string(fn), Err: e})
	}
	// Return nil
	return nil
}

// EvalGoldenFile evaluates the testcase if it equals the test data from the golden file provided by the testcase name.
// It returns an error if the testcase data does not equal the contents of the golden file. The golden file must reside in the default golden files directory
// testdata/ with the default golden file type .golden.
func EvalGoldenFile(tc *Testcase) error {
	// Return an error if tc is nil
	if tc == nil {
		return tserr.NilPtr()
	}
	// Retrieve golden file path for testcase name
	fn, e := GoldenFilePath(tc.Name)
	// Return an error if goldenPath fails
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "goldenPath", Fn: tc.Name, Err: e})
	}
	// Retrieve the reference data from golden file provided by the testcase name
	ref, e := ReadFile(fn)
	// Return an error if ReadFile fails
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(fn), Err: e})
	}
	// Normalize new lines in ref
	refn := NormNewlinesStr(string(ref))
	// Normalize new lines in test data
	test := NormNewlinesStr(tc.Data)
	// Return an error if the testcase data does not equal the contents of the golden file
	if test != refn {
		return tserr.EqualStr(&tserr.EqualStrArgs{Var: tc.Name, Actual: tc.Data, Want: string(ref)})
	}
	// Return nil
	return nil
}
