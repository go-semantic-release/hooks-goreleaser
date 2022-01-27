// Package publish contains the publishing pipe.
package publish

import (
	"fmt"

	"github.com/goreleaser/goreleaser/int/middleware/errhandler"
	"github.com/goreleaser/goreleaser/int/middleware/logging"
	"github.com/goreleaser/goreleaser/int/middleware/skip"
	"github.com/goreleaser/goreleaser/int/pipe/artifactory"
	"github.com/goreleaser/goreleaser/int/pipe/aur"
	"github.com/goreleaser/goreleaser/int/pipe/blob"
	"github.com/goreleaser/goreleaser/int/pipe/brew"
	"github.com/goreleaser/goreleaser/int/pipe/custompublishers"
	"github.com/goreleaser/goreleaser/int/pipe/docker"
	"github.com/goreleaser/goreleaser/int/pipe/gofish"
	"github.com/goreleaser/goreleaser/int/pipe/krew"
	"github.com/goreleaser/goreleaser/int/pipe/milestone"
	"github.com/goreleaser/goreleaser/int/pipe/release"
	"github.com/goreleaser/goreleaser/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/int/pipe/sign"
	"github.com/goreleaser/goreleaser/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/int/pipe/upload"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Publisher should be implemented by pipes that want to publish artifacts.
type Publisher interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Publish(ctx *context.Context) error
}

// nolint: gochecknoglobals
var publishers = []Publisher{
	blob.Pipe{},
	upload.Pipe{},
	custompublishers.Pipe{},
	artifactory.Pipe{},
	docker.Pipe{},
	docker.ManifestPipe{},
	sign.DockerPipe{},
	snapcraft.Pipe{},
	// This should be one of the last steps
	release.Pipe{},
	// brew et al use the release URL, so, they should be last
	brew.Pipe{},
	aur.Pipe{},
	gofish.Pipe{},
	krew.Pipe{},
	scoop.Pipe{},
	milestone.Pipe{},
}

// Pipe that publishes artifacts.
type Pipe struct{}

func (Pipe) String() string                 { return "publishing" }
func (Pipe) Skip(ctx *context.Context) bool { return ctx.SkipPublish }

func (Pipe) Run(ctx *context.Context) error {
	for _, publisher := range publishers {
		if err := skip.Maybe(
			publisher,
			logging.Log(
				publisher.String(),
				errhandler.Handle(publisher.Publish),
				logging.ExtraPadding,
			),
		)(ctx); err != nil {
			return fmt.Errorf("%s: failed to publish artifacts: %w", publisher.String(), err)
		}
	}
	return nil
}
