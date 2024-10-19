// Package defaults make the list of Defaulter implementations available
// so projects extending GoReleaser are able to use it, namely, GoDownloader.
package defaults

import (
	"fmt"

	"github.com/goreleaser/goreleaser/v2/int/pipe/archive"
	"github.com/goreleaser/goreleaser/v2/int/pipe/artifactory"
	"github.com/goreleaser/goreleaser/v2/int/pipe/aur"
	"github.com/goreleaser/goreleaser/v2/int/pipe/blob"
	"github.com/goreleaser/goreleaser/v2/int/pipe/bluesky"
	"github.com/goreleaser/goreleaser/v2/int/pipe/brew"
	"github.com/goreleaser/goreleaser/v2/int/pipe/build"
	"github.com/goreleaser/goreleaser/v2/int/pipe/changelog"
	"github.com/goreleaser/goreleaser/v2/int/pipe/checksums"
	"github.com/goreleaser/goreleaser/v2/int/pipe/chocolatey"
	"github.com/goreleaser/goreleaser/v2/int/pipe/discord"
	"github.com/goreleaser/goreleaser/v2/int/pipe/dist"
	"github.com/goreleaser/goreleaser/v2/int/pipe/docker"
	"github.com/goreleaser/goreleaser/v2/int/pipe/gomod"
	"github.com/goreleaser/goreleaser/v2/int/pipe/ko"
	"github.com/goreleaser/goreleaser/v2/int/pipe/krew"
	"github.com/goreleaser/goreleaser/v2/int/pipe/linkedin"
	"github.com/goreleaser/goreleaser/v2/int/pipe/mastodon"
	"github.com/goreleaser/goreleaser/v2/int/pipe/mattermost"
	"github.com/goreleaser/goreleaser/v2/int/pipe/milestone"
	"github.com/goreleaser/goreleaser/v2/int/pipe/nfpm"
	"github.com/goreleaser/goreleaser/v2/int/pipe/nix"
	"github.com/goreleaser/goreleaser/v2/int/pipe/notary"
	"github.com/goreleaser/goreleaser/v2/int/pipe/opencollective"
	"github.com/goreleaser/goreleaser/v2/int/pipe/project"
	"github.com/goreleaser/goreleaser/v2/int/pipe/reddit"
	"github.com/goreleaser/goreleaser/v2/int/pipe/release"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sbom"
	"github.com/goreleaser/goreleaser/v2/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sign"
	"github.com/goreleaser/goreleaser/v2/int/pipe/slack"
	"github.com/goreleaser/goreleaser/v2/int/pipe/smtp"
	"github.com/goreleaser/goreleaser/v2/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/v2/int/pipe/snapshot"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sourcearchive"
	"github.com/goreleaser/goreleaser/v2/int/pipe/teams"
	"github.com/goreleaser/goreleaser/v2/int/pipe/telegram"
	"github.com/goreleaser/goreleaser/v2/int/pipe/twitter"
	"github.com/goreleaser/goreleaser/v2/int/pipe/universalbinary"
	"github.com/goreleaser/goreleaser/v2/int/pipe/upload"
	"github.com/goreleaser/goreleaser/v2/int/pipe/upx"
	"github.com/goreleaser/goreleaser/v2/int/pipe/webhook"
	"github.com/goreleaser/goreleaser/v2/int/pipe/winget"
	"github.com/goreleaser/goreleaser/v2/pkg/context"
)

// Defaulter can be implemented by a Piper to set default values for its
// configuration.
type Defaulter interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Default(ctx *context.Context) error
}

// Defaulters is the list of defaulters.
//
//nolint:gochecknoglobals
var Defaulters = []Defaulter{
	dist.Pipe{},
	snapshot.Pipe{},
	release.Pipe{},
	project.Pipe{},
	changelog.Pipe{},
	gomod.Pipe{},
	build.Pipe{},
	universalbinary.Pipe{},
	sign.BinaryPipe{},
	notary.MacOS{},
	upx.Pipe{},
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
	upload.Pipe{},
	aur.Pipe{},
	nix.Pipe{},
	winget.Pipe{},
	brew.Pipe{},
	krew.Pipe{},
	ko.Pipe{},
	scoop.Pipe{},
	discord.Pipe{},
	reddit.Pipe{},
	slack.Pipe{},
	teams.Pipe{},
	twitter.Pipe{},
	smtp.Pipe{},
	mastodon.Pipe{},
	mattermost.Pipe{},
	milestone.Pipe{},
	linkedin.Pipe{},
	telegram.Pipe{},
	webhook.Pipe{},
	chocolatey.Pipe{},
	opencollective.Pipe{},
	bluesky.Pipe{},
}
