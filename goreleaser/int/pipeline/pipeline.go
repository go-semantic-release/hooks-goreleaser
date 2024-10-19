// Package pipeline provides generic errors for pipes to use.
package pipeline

import (
	"fmt"

	"github.com/goreleaser/goreleaser/v2/int/pipe/announce"
	"github.com/goreleaser/goreleaser/v2/int/pipe/archive"
	"github.com/goreleaser/goreleaser/v2/int/pipe/aur"
	"github.com/goreleaser/goreleaser/v2/int/pipe/before"
	"github.com/goreleaser/goreleaser/v2/int/pipe/brew"
	"github.com/goreleaser/goreleaser/v2/int/pipe/build"
	"github.com/goreleaser/goreleaser/v2/int/pipe/changelog"
	"github.com/goreleaser/goreleaser/v2/int/pipe/checksums"
	"github.com/goreleaser/goreleaser/v2/int/pipe/chocolatey"
	"github.com/goreleaser/goreleaser/v2/int/pipe/defaults"
	"github.com/goreleaser/goreleaser/v2/int/pipe/dist"
	"github.com/goreleaser/goreleaser/v2/int/pipe/docker"
	"github.com/goreleaser/goreleaser/v2/int/pipe/effectiveconfig"
	"github.com/goreleaser/goreleaser/v2/int/pipe/env"
	"github.com/goreleaser/goreleaser/v2/int/pipe/git"
	"github.com/goreleaser/goreleaser/v2/int/pipe/gomod"
	"github.com/goreleaser/goreleaser/v2/int/pipe/ko"
	"github.com/goreleaser/goreleaser/v2/int/pipe/krew"
	"github.com/goreleaser/goreleaser/v2/int/pipe/metadata"
	"github.com/goreleaser/goreleaser/v2/int/pipe/nfpm"
	"github.com/goreleaser/goreleaser/v2/int/pipe/nix"
	"github.com/goreleaser/goreleaser/v2/int/pipe/notary"
	"github.com/goreleaser/goreleaser/v2/int/pipe/partial"
	"github.com/goreleaser/goreleaser/v2/int/pipe/prebuild"
	"github.com/goreleaser/goreleaser/v2/int/pipe/publish"
	"github.com/goreleaser/goreleaser/v2/int/pipe/reportsizes"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sbom"
	"github.com/goreleaser/goreleaser/v2/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/v2/int/pipe/semver"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sign"
	"github.com/goreleaser/goreleaser/v2/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/v2/int/pipe/snapshot"
	"github.com/goreleaser/goreleaser/v2/int/pipe/sourcearchive"
	"github.com/goreleaser/goreleaser/v2/int/pipe/universalbinary"
	"github.com/goreleaser/goreleaser/v2/int/pipe/upx"
	"github.com/goreleaser/goreleaser/v2/int/pipe/winget"
	"github.com/goreleaser/goreleaser/v2/pkg/context"
)

// Piper defines a pipe, which can be part of a pipeline (a series of pipes).
type Piper interface {
	fmt.Stringer

	// Run the pipe
	Run(ctx *context.Context) error
}

// BuildPipeline contains all build-related pipe implementations in order.
//
//nolint:gochecknoglobals
var BuildPipeline = []Piper{
	// set default dist folder and remove it if `--clean` is set
	dist.CleanPipe{},
	// load and validate environment variables
	env.Pipe{},
	// get and validate git repo state
	git.Pipe{},
	// parse current tag to a semver
	semver.Pipe{},
	// load default configs
	defaults.Pipe{},
	// setup things for partial builds/releases
	partial.Pipe{},
	// snapshot version handling
	snapshot.Pipe{},
	// run global hooks before build
	before.Pipe{},
	// ensure ./dist exists and is empty
	dist.Pipe{},
	// setup metadata options
	metadata.Pipe{},
	// creates a metadta.json files in the dist directory
	metadata.MetaPipe{},
	// setup gomod-related stuff
	gomod.Pipe{},
	// run prebuild stuff
	prebuild.Pipe{},
	// proxy gomod if needed
	gomod.CheckGoModPipe{},
	// proxy gomod if needed
	gomod.ProxyPipe{},
	// writes the actual config (with defaults et al set) to dist
	effectiveconfig.Pipe{},
	// build
	build.Pipe{},
	// universal binary handling
	universalbinary.Pipe{},
	// sign binaries
	sign.BinaryPipe{},
	// notarize macos apps
	notary.MacOS{},
	// upx
	upx.Pipe{},
}

// BuildCmdPipeline is the pipeline run by goreleaser build.
//
//nolint:gochecknoglobals
var BuildCmdPipeline = append(
	BuildPipeline,
	reportsizes.Pipe{},
	metadata.ArtifactsPipe{},
)

// Pipeline contains all pipe implementations in order.
//
//nolint:gochecknoglobals
var Pipeline = append(
	BuildPipeline,
	// builds the release changelog
	changelog.Pipe{},
	// archive in tar.gz, zip or binary (which does no archiving at all)
	archive.Pipe{},
	// archive the source code using git-archive
	sourcearchive.Pipe{},
	// archive via fpm (deb, rpm) using "native" go impl
	nfpm.Pipe{},
	// archive via snapcraft (snap)
	snapcraft.Pipe{},
	// create SBOMs of artifacts
	sbom.Pipe{},
	// checksums of the files
	checksums.Pipe{},
	// sign artifacts
	sign.Pipe{},
	// create arch linux aur pkgbuild
	aur.Pipe{},
	// create nixpkgs
	nix.NewBuild(),
	// winget installers
	winget.Pipe{},
	// create brew tap
	brew.Pipe{},
	// krew plugins
	krew.Pipe{},
	// create scoop buckets
	scoop.Pipe{},
	// create chocolatey pkg and publish
	chocolatey.Pipe{},
	// reports artifacts sizes to the log and to artifacts.json
	reportsizes.Pipe{},
	// create and push docker images
	docker.Pipe{},
	// create and push docker images using ko
	ko.Pipe{},
	// publishes artifacts
	publish.New(),
	// creates a artifacts.json files in the dist directory
	metadata.ArtifactsPipe{},
	// announce releases
	announce.Pipe{},
)
