package hooks

import (
	"log"
	"os"

	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
)

type GoReleaser struct {
	logger       *log.Logger
	providerName string
	prerelease   bool
}

func (t *GoReleaser) Init(m map[string]string) error {
	t.logger = log.New(os.Stderr, "", 0)
	t.logger.Printf("init: %v\n", m)
	return nil
}

func (t *GoReleaser) Name() string {
	return "GoReleaser"
}

func (t *GoReleaser) Version() string {
	return "dev"
}

func (t *GoReleaser) Success(config *hooks.SuccessHookConfig) error {
	t.logger.Println("old version: " + config.PrevRelease.Version)
	t.logger.Println("new version: " + config.NewRelease.Version)
	return nil
}

func (t *GoReleaser) NoRelease(config *hooks.NoReleaseConfig) error {
	return nil
}
