package hooks

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/caarlos0/log"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/goreleaser/goreleaser/int/middleware/errhandler"
	"github.com/goreleaser/goreleaser/int/middleware/logging"
	"github.com/goreleaser/goreleaser/int/middleware/skip"
	"github.com/goreleaser/goreleaser/int/pipe/git"
	"github.com/goreleaser/goreleaser/int/pipeline"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

var HVERSION = "dev"

type GoReleaser struct {
	providerName  string
	currentBranch string
	prerelease    bool
	parallelism   int
}

func (gr *GoReleaser) Init(m map[string]string) error {
	gr.providerName = m["provider"]
	gr.currentBranch = m["currentBranch"]
	gr.prerelease = m["prerelease"] == "true"
	gr.parallelism = runtime.NumCPU()
	if m["parallelism"] == "" {
		return nil
	}
	val, err := strconv.ParseUint(m["parallelism"], 10, 32)
	if err != nil {
		log.Warnf("could not parse parallelism=%s", m["parallelism"])
		return nil
	}
	gr.parallelism = int(val)
	return nil
}

func (gr *GoReleaser) Success(shConfig *hooks.SuccessHookConfig) error {
	newVersion := shConfig.NewRelease.Version
	currentSha := shConfig.NewRelease.SHA
	cfg, err := readGoReleaserConfig()
	if err != nil {
		return err
	}
	ctx, cancel := context.NewWithTimeout(cfg, 30*time.Minute)
	defer cancel()

	// set some fixed parameters for GoReleaser
	ctx.Parallelism = gr.parallelism
	log.Infof("build parallelism is set to %d", ctx.Parallelism)
	ctx.Clean = true
	ctx.Config.Changelog.Skip = "true" // never generate changelog

	ctx.Version = newVersion
	ctx.Git = context.GitInfo{
		Branch:      gr.currentBranch,
		CurrentTag:  fmt.Sprintf("v%s", newVersion),
		Commit:      currentSha,
		FullCommit:  currentSha,
		ShortCommit: currentSha,
	}

	repo := config.Repo{
		Owner: shConfig.RepoInfo.Owner,
		Name:  shConfig.RepoInfo.Repo,
	}
	switch gr.providerName {
	case "GitLab":
		if ctx.Config.Release.GitLab.Name == "" {
			ctx.Config.Release.GitLab = repo
		}
	default:
		if ctx.Config.Release.GitHub.Name == "" {
			ctx.Config.Release.GitHub = repo
		}
	}

	// keep release config in sync with the already existing release from semantic-release
	ctx.Config.Release.ReleaseNotesMode = config.ReleaseNotesModeKeepExisting
	ctx.Config.Release.NameTemplate = ctx.Git.CurrentTag
	ctx.PreRelease = gr.prerelease
	ctx.Config.Release.Draft = false // always disable drafts

	for _, pipe := range pipeline.Pipeline {
		if _, ok := pipe.(git.Pipe); ok {
			log.Info("skipping git pipe")
			continue
		}
		err := skip.Maybe(pipe, logging.Log(pipe.String(), errhandler.Handle(pipe.Run)))(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gr *GoReleaser) NoRelease(_ *hooks.NoReleaseConfig) error {
	return nil
}

func (gr *GoReleaser) Name() string {
	return "GoReleaser"
}

func (gr *GoReleaser) Version() string {
	return HVERSION
}
