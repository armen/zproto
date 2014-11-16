// +build test_zproto_c
package example

import (
	"testing"
)

// Tests interoperability of Go implementation with C implementation
func TestLogZprotoC(t *testing.T) {
	testLogInputFromC(t)
}
