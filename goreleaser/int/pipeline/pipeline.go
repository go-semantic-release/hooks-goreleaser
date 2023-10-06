// Package pipeline provides generic errors for pipes to use.
package pipeline

import (
	"fmt"

	"github.com/goreleaser/goreleaser/int/pipe/announce"
	"github.com/goreleaser/goreleaser/int/pipe/archive"
	"github.com/goreleaser/goreleaser/int/pipe/aur"
	"github.com/goreleaser/goreleaser/int/pipe/before"
	"github.com/goreleaser/goreleaser/int/pipe/brew"
	"github.com/goreleaser/goreleaser/int/pipe/build"
	"github.com/goreleaser/goreleaser/int/pipe/changelog"
	"github.com/goreleaser/goreleaser/int/pipe/checksums"
	"github.com/goreleaser/goreleaser/int/pipe/chocolatey"
	"github.com/goreleaser/goreleaser/int/pipe/defaults"
	"github.com/goreleaser/goreleaser/int/pipe/dist"
	"github.com/goreleaser/goreleaser/int/pipe/docker"
	"github.com/goreleaser/goreleaser/int/pipe/effectiveconfig"
	"github.com/goreleaser/goreleaser/int/pipe/env"
	"github.com/goreleaser/goreleaser/int/pipe/git"
	"github.com/goreleaser/goreleaser/int/pipe/gomod"
	"github.com/goreleaser/goreleaser/int/pipe/krew"
	"github.com/goreleaser/goreleaser/int/pipe/metadata"
	"github.com/goreleaser/goreleaser/int/pipe/nfpm"
	"github.com/goreleaser/goreleaser/int/pipe/nix"
	"github.com/goreleaser/goreleaser/int/pipe/prebuild"
	"github.com/goreleaser/goreleaser/int/pipe/publish"
	"github.com/goreleaser/goreleaser/int/pipe/reportsizes"
	"github.com/goreleaser/goreleaser/int/pipe/sbom"
	"github.com/goreleaser/goreleaser/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/int/pipe/semver"
	"github.com/goreleaser/goreleaser/int/pipe/sign"
	"github.com/goreleaser/goreleaser/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/int/pipe/snapshot"
	"github.com/goreleaser/goreleaser/int/pipe/sourcearchive"
	"github.com/goreleaser/goreleaser/int/pipe/universalbinary"
	"github.com/goreleaser/goreleaser/int/pipe/upx"
	"github.com/goreleaser/goreleaser/int/pipe/winget"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Piper defines a pipe, which can be part of a pipeline (a series of pipes).
type Piper interface {
	fmt.Stringer

	// Run the pipe
	Run(ctx *context.Context) error
}

// BuildPipeline contains all build-related pipe implementations in order.
// nolint:gochecknoglobals
var BuildPipeline = []Piper{
	// load and validate environment variables
	env.Pipe{},
	// get and validate git repo state
	git.Pipe{},
	// parse current tag to a semver
	semver.Pipe{},
	// load default configs
	defaults.Pipe{},
	// snapshot version handling
	snapshot.Pipe{},
	// run global hooks before build
	before.Pipe{},
	// ensure ./dist is clean
	dist.Pipe{},
	// setup gomod-related stuff
	gomod.Pipe{},
	// run prebuild stuff
	prebuild.Pipe{},
	// proxy gomod if needed
	gomod.ProxyPipe{},
	// writes the actual config (with defaults et al set) to dist
	effectiveconfig.Pipe{},
	// build
	build.Pipe{},
	// universal binary handling
	universalbinary.Pipe{},
	// upx
	upx.Pipe{},
}

// BuildCmdPipeline is the pipeline run by goreleaser build.
// nolint:gochecknoglobals
var BuildCmdPipeline = append(
	BuildPipeline,
	reportsizes.Pipe{},
	metadata.Pipe{},
)

// Pipeline contains all pipe implementations in order.
// nolint: gochecknoglobals
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
	// publishes artifacts
	publish.New(),
	// creates a metadata.json and an artifacts.json files in the dist folder
	metadata.Pipe{},
	// announce releases
	announce.Pipe{},
)
