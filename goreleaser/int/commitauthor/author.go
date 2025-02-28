// Package commitauthor provides common commit author functionality.
package commitauthor

import (
	"github.com/goreleaser/goreleaser/v2/int/tmpl"
	"github.com/goreleaser/goreleaser/v2/pkg/config"
	"github.com/goreleaser/goreleaser/v2/pkg/context"
)

const (
	defaultName  = "goreleaserbot"
	defaultEmail = "bot@goreleaser.com"
)

// Get templates the commit author and returns the filled fields.
func Get(ctx *context.Context, og config.CommitAuthor) (config.CommitAuthor, error) {
	var author config.CommitAuthor
	var err error

	author.Name, err = tmpl.New(ctx).Apply(og.Name)
	if err != nil {
		return author, err
	}
	author.Email, err = tmpl.New(ctx).Apply(og.Email)
	return author, err
}

// Default sets the default commit author name and email.
func Default(og config.CommitAuthor) config.CommitAuthor {
	if og.Name == "" {
		og.Name = defaultName
	}
	if og.Email == "" {
		og.Email = defaultEmail
	}
	return og
}
