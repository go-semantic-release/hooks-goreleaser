package testlib

import (
	"testing"

	"github.com/goreleaser/goreleaser/int/pipe"
)

func TestAssertSkipped(t *testing.T) {
	AssertSkipped(t, pipe.Skip("skip"))
}
