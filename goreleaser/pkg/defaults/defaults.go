// Package defaults make the list of Defaulter implementations available
// so projects extending GoReleaser are able to use it, namely, GoDownloader.
package defaults

import (
	"fmt"

	"github.com/goreleaser/goreleaser/int/pipe/archive"
	"github.com/goreleaser/goreleaser/int/pipe/artifactory"
	"github.com/goreleaser/goreleaser/int/pipe/aur"
	"github.com/goreleaser/goreleaser/int/pipe/blob"
	"github.com/goreleaser/goreleaser/int/pipe/brew"
	"github.com/goreleaser/goreleaser/int/pipe/build"
	"github.com/goreleaser/goreleaser/int/pipe/checksums"
	"github.com/goreleaser/goreleaser/int/pipe/discord"
	"github.com/goreleaser/goreleaser/int/pipe/docker"
	"github.com/goreleaser/goreleaser/int/pipe/gofish"
	"github.com/goreleaser/goreleaser/int/pipe/gomod"
	"github.com/goreleaser/goreleaser/int/pipe/krew"
	"github.com/goreleaser/goreleaser/int/pipe/linkedin"
	"github.com/goreleaser/goreleaser/int/pipe/mattermost"
	"github.com/goreleaser/goreleaser/int/pipe/milestone"
	"github.com/goreleaser/goreleaser/int/pipe/nfpm"
	"github.com/goreleaser/goreleaser/int/pipe/project"
	"github.com/goreleaser/goreleaser/int/pipe/reddit"
	"github.com/goreleaser/goreleaser/int/pipe/release"
	"github.com/goreleaser/goreleaser/int/pipe/sbom"
	"github.com/goreleaser/goreleaser/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/int/pipe/sign"
	"github.com/goreleaser/goreleaser/int/pipe/slack"
	"github.com/goreleaser/goreleaser/int/pipe/smtp"
	"github.com/goreleaser/goreleaser/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/int/pipe/snapshot"
	"github.com/goreleaser/goreleaser/int/pipe/sourcearchive"
	"github.com/goreleaser/goreleaser/int/pipe/teams"
	"github.com/goreleaser/goreleaser/int/pipe/telegram"
	"github.com/goreleaser/goreleaser/int/pipe/twitter"
	"github.com/goreleaser/goreleaser/int/pipe/universalbinary"
	"github.com/goreleaser/goreleaser/int/pipe/webhook"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Defaulter can be implemented by a Piper to set default values for its
// configuration.
type Defaulter interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Default(ctx *context.Context) error
}

// Defaulters is the list of defaulters.
// nolint: gochecknoglobals
var Defaulters = []Defaulter{
	snapshot.Pipe{},
	release.Pipe{},
	project.Pipe{},
	gomod.Pipe{},
	build.Pipe{},
	universalbinary.Pipe{},
	sourcearchive.Pipe{},
	archive.Pipe{},
	nfpm.Pipe{},
	snapcraft.Pipe{},
	checksums.Pipe{},
	sign.Pipe{},
	sign.DockerPipe{},
	sbom.Pipe{},
	docker.Pipe{},
	docker.ManifestPipe{},
	artifactory.Pipe{},
	blob.Pipe{},
	aur.Pipe{},
	brew.Pipe{},
	krew.Pipe{},
	gofish.Pipe{},
	scoop.Pipe{},
	discord.Pipe{},
	reddit.Pipe{},
	slack.Pipe{},
	teams.Pipe{},
	twitter.Pipe{},
	smtp.Pipe{},
	mattermost.Pipe{},
	milestone.Pipe{},
	linkedin.Pipe{},
	telegram.Pipe{},
	webhook.Pipe{},
}
