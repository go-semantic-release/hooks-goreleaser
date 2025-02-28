// Package publish contains the publishing pipe.
package publish

import (
	"fmt"

	"github.com/goreleaser/goreleaser/v2/int/middleware/errhandler"
	"github.com/goreleaser/goreleaser/v2/int/middleware/logging"
	"github.com/goreleaser/goreleaser/v2/int/middleware/skip"
	"github.com/goreleaser/goreleaser/v2/int/pipe/artifactory"
	"github.com/goreleaser/goreleaser/v2/int/pipe/aur"
	"github.com/goreleaser/goreleaser/v2/int/pipe/blob"
	"github.com/goreleaser/goreleaser/v2/int/pipe/brew"
	"github.com/goreleaser/goreleaser/v2/int/pipe/chocolatey"
	"github.com/goreleaser/goreleaser/v2/int/pipe/custompublishers"
	"github.com/goreleaser/goreleaser/v2/int/pipe/docker"
	"github.com/goreleaser/goreleaser/v2/int/pipe/ko"
	"github.com/goreleaser/goreleaser/v2/int/pipe/krew"
	"github.com/goreleaser/goreleaser/v2/int/pipe/milestone"
	"github.com/goreleaser/goreleaser/v2/int/pipe/nix"
	"github.com/goreleaser/goreleaser/v2/int/pipe/release"
	"github.com/goreleaser/goreleaser/v2/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sign"
	"github.com/goreleaser/goreleaser/v2/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/v2/int/pipe/upload"
	"github.com/goreleaser/goreleaser/v2/int/pipe/winget"
	"github.com/goreleaser/goreleaser/v2/int/skips"
	"github.com/goreleaser/goreleaser/v2/pkg/context"
)

// Publisher should be implemented by pipes that want to publish artifacts.
type Publisher interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Publish(ctx *context.Context) error
}

// New publish pipeline.
func New() Pipe {
	return Pipe{
		pipeline: []Publisher{
			blob.Pipe{},
			upload.Pipe{},
			artifactory.Pipe{},
			custompublishers.Pipe{},
			docker.Pipe{},
			docker.ManifestPipe{},
			ko.Pipe{},
			sign.DockerPipe{},
			snapcraft.Pipe{},
			// This should be one of the last steps
			release.Pipe{},
			// brew et al use the release URL, so, they should be last
			nix.NewPublish(),
			winget.Pipe{},
			brew.Pipe{},
			aur.Pipe{},
			krew.Pipe{},
			scoop.Pipe{},
			chocolatey.Pipe{},
			milestone.Pipe{},
		},
	}
}

// Pipe that publishes artifacts.
type Pipe struct {
	pipeline []Publisher
}

func (Pipe) String() string                 { return "publishing" }
func (Pipe) Skip(ctx *context.Context) bool { return skips.Any(ctx, skips.Publish) }

func (p Pipe) Run(ctx *context.Context) error {
	memo := errhandler.Memo{}
	for _, publisher := range p.pipeline {
		if err := skip.Maybe(
			publisher,
			logging.PadLog(
				publisher.String(),
				errhandler.Handle(publisher.Publish),
			),
		)(ctx); err != nil {
			if ig, ok := publisher.(Continuable); ok && ig.ContinueOnError() && !ctx.FailFast {
				memo.Memorize(fmt.Errorf("%s: %w", publisher.String(), err))
				continue
			}
			return fmt.Errorf("%s: failed to publish artifacts: %w", publisher.String(), err)
		}
	}
	return memo.Error()
}

type Continuable interface {
	ContinueOnError() bool
}
