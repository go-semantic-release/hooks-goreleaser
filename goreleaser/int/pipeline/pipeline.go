// Package pipeline provides generic erros for pipes to use.
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
	"github.com/goreleaser/goreleaser/int/pipe/defaults"
	"github.com/goreleaser/goreleaser/int/pipe/dist"
	"github.com/goreleaser/goreleaser/int/pipe/docker"
	"github.com/goreleaser/goreleaser/int/pipe/effectiveconfig"
	"github.com/goreleaser/goreleaser/int/pipe/env"
	"github.com/goreleaser/goreleaser/int/pipe/git"
	"github.com/goreleaser/goreleaser/int/pipe/gofish"
	"github.com/goreleaser/goreleaser/int/pipe/gomod"
	"github.com/goreleaser/goreleaser/int/pipe/krew"
	"github.com/goreleaser/goreleaser/int/pipe/metadata"
	"github.com/goreleaser/goreleaser/int/pipe/nfpm"
	"github.com/goreleaser/goreleaser/int/pipe/prebuild"
	"github.com/goreleaser/goreleaser/int/pipe/publish"
	"github.com/goreleaser/goreleaser/int/pipe/sbom"
	"github.com/goreleaser/goreleaser/int/pipe/scoop"
	"github.com/goreleaser/goreleaser/int/pipe/semver"
	"github.com/goreleaser/goreleaser/int/pipe/sign"
	"github.com/goreleaser/goreleaser/int/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/int/pipe/snapshot"
	"github.com/goreleaser/goreleaser/int/pipe/sourcearchive"
	"github.com/goreleaser/goreleaser/int/pipe/universalbinary"
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
	env.Pipe{},             // load and validate environment variables
	git.Pipe{},             // get and validate git repo state
	semver.Pipe{},          // parse current tag to a semver
	before.Pipe{},          // run global hooks before build
	defaults.Pipe{},        // load default configs
	snapshot.Pipe{},        // snapshot version handling
	dist.Pipe{},            // ensure ./dist is clean
	gomod.Pipe{},           // setup gomod-related stuff
	prebuild.Pipe{},        // run prebuild stuff
	gomod.ProxyPipe{},      // proxy gomod if needed
	effectiveconfig.Pipe{}, // writes the actual config (with defaults et al set) to dist
	changelog.Pipe{},       // builds the release changelog
	build.Pipe{},           // build
	universalbinary.Pipe{}, // universal binary handling
}

// BuildCmdPipeline is the pipeline run by goreleaser build.
// nolint:gochecknoglobals
var BuildCmdPipeline = append(BuildPipeline, metadata.Pipe{})

// Pipeline contains all pipe implementations in order.
// nolint: gochecknoglobals
var Pipeline = append(
	BuildPipeline,
	archive.Pipe{},       // archive in tar.gz, zip or binary (which does no archiving at all)
	sourcearchive.Pipe{}, // archive the source code using git-archive
	nfpm.Pipe{},          // archive via fpm (deb, rpm) using "native" go impl
	snapcraft.Pipe{},     // archive via snapcraft (snap)
	aur.Pipe{},           // create arch linux aur pkgbuild
	brew.Pipe{},          // create brew tap
	gofish.Pipe{},        // create gofish rig
	krew.Pipe{},          // krew plugins
	scoop.Pipe{},         // create scoop buckets
	sbom.Pipe{},          // create SBOMs of artifacts
	checksums.Pipe{},     // checksums of the files
	sign.Pipe{},          // sign artifacts
	docker.Pipe{},        // create and push docker images
	metadata.Pipe{},      // creates a metadata.json and an artifacts.json files in the dist folder
	publish.Pipe{},       // publishes artifacts
	announce.Pipe{},      // announce releases
)
