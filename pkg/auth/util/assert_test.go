package util_test

import (
	"testing"

	"github.com/stanekondrej/quarkchess/auth/pkg/auth/util"
)

func TestAssertion(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	util.Assert(false)
}
