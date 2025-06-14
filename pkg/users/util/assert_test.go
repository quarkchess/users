package util_test

import (
	"testing"

	"github.com/quarkchess/users/pkg/users/util"
)

func TestAssertion(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	util.Assert(false)
}
