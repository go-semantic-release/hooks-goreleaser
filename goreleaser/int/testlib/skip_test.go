package testlib

import (
	"testing"

	"github.com/goreleaser/goreleaser/v2/int/pipe"
)

func TestAssertSkipped(t *testing.T) {
	AssertSkipped(t, pipe.Skip("skip"))
}
