// Copyright (c) 2023-2026 thorsphere.
// All Rights Reserved. Use is governed by the Functional Source License v1.1
// (FSL-1.1-ALv2) that can be found in the LICENSE file.
package tsfio

// Import go standard packages and tserr
import (
	"bytes"   // bytes
	"strings" // strings

	"github.com/thorsphere/tserr" // tserr
)

// NormNewlinesBytes normalizes new lines in the byte slice i. New lines are normalized to the Unix representation of a new line as line feed LF (0x0A).
// Therefore, Windows new lines CR LF (0x0D 0x0A) are replaced by Unix new lines LF (0x0A). Also, Mac new lines CR (0x0D) are replaced by Unix new lines LF (0x0A).
// It returns a normalized copy of i. If i is nil, it returns an error.
func NormNewlinesBytes(i []byte) ([]byte, error) {
	// Return nil and an error, if i is nil
	if i == nil {
		return nil, tserr.NilPtr()
	}
	// Windows new lines CR LF (0x0D 0x0A) are replaced by Unix new lines LF (0x0A).
	i = bytes.ReplaceAll(i, []byte("\r\n"), []byte("\n"))
	// Mac new lines CR (0x0D) are replaced by Unix new lines LF (0x0A).
	i = bytes.ReplaceAll(i, []byte("\r"), []byte("\n")) // MAC
	// Return the normalized copy of i.
	return i, nil
}

// NormNewlinesStr normalizes new lines in the string i. New lines are normalized to the Unix representation of a new line as line feed LF (0x0A).
// Therefore, Windows new lines CR LF (0x0D 0x0A) are replaced by Unix new lines LF (0x0A). Also, Mac new lines CR (0x0D) are replaced by Unix new lines LF (0x0A).
// It returns a normalized copy of i.
func NormNewlinesStr(i string) string {
	// Windows new lines CR LF (0x0D 0x0A) are replaced by Unix new lines LF (0x0A).
	i = strings.ReplaceAll(i, "\r\n", "\n")
	// Mac new lines CR (0x0D) are replaced by Unix new lines LF (0x0A).
	i = strings.ReplaceAll(i, "\r", "\n")
	// Return the normalized copy of i.
	return i
}
