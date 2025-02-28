// Package healthcheck checks for missing binaries that the user needs to
// install.
package healthcheck

import (
	"fmt"

	"github.com/goreleaser/goreleaser/v2/int/pipe/chocolatey"
	"github.com/goreleaser/goreleaser/v2/int/pipe/docker"
	"github.com/goreleaser/goreleaser/v2/int/pipe/nix"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sbom"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sign"
	"github.com/goreleaser/goreleaser/v2/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/v2/pkg/context"
)

// Healthchecker should be implemented by pipes that want checks.
type Healthchecker interface {
	fmt.Stringer

	// Dependencies return the binaries of the dependencies needed.
	Dependencies(ctx *context.Context) []string
}

// Healthcheckers is the list of healthchekers.
//
//nolint:gochecknoglobals
var Healthcheckers = []Healthchecker{
	system{},
	snapcraft.Pipe{},
	sign.Pipe{},
	sign.BinaryPipe{},
	sign.DockerPipe{},
	sbom.Pipe{},
	docker.Pipe{},
	docker.ManifestPipe{},
	chocolatey.Pipe{},
	nix.NewPublish(),
}

type system struct{}

func (system) String() string                           { return "system" }
func (system) Dependencies(_ *context.Context) []string { return []string{"git", "go"} }
