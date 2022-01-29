package main

import (
	"github.com/apex/log"
	"github.com/fatih/color"
	hooksGoReleaser "github.com/go-semantic-release/hooks-goreleaser/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

func main() {
	log.SetHandler(hooksGoReleaser.NewLogHandler())
	color.NoColor = true
	plugin.Serve(&plugin.ServeOpts{
		Hooks: func() hooks.Hooks {
			return &hooksGoReleaser.GoReleaser{}
		},
	})
}
