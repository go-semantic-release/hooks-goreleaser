package main

import (
	hooksGoReleaser "github.com/go-semantic-release/hooks-goreleaser/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Hooks: func() hooks.Hooks {
			return &hooksGoReleaser.GoReleaser{}
		},
	})
}
