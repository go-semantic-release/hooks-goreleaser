// Package snapshot provides the snapshotting functionality to goreleaser.
package snapshot

import (
	"fmt"

	"github.com/caarlos0/log"
	"github.com/goreleaser/goreleaser/v2/int/deprecate"
	"github.com/goreleaser/goreleaser/v2/int/tmpl"
	"github.com/goreleaser/goreleaser/v2/pkg/context"
)

// Pipe for setting up the snapshot feature..
type Pipe struct{}

func (Pipe) String() string                 { return "snapshotting" }
func (Pipe) Skip(ctx *context.Context) bool { return !ctx.Snapshot }

// Default sets the pipe defaults.
func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Snapshot.VersionTemplate == "" {
		ctx.Config.Snapshot.VersionTemplate = "{{ .Version }}-SNAPSHOT-{{ .ShortCommit }}"
	}
	if ctx.Config.Snapshot.NameTemplate != "" {
		deprecate.Notice(ctx, "snapshot.name_template")
		ctx.Config.Snapshot.VersionTemplate = ctx.Config.Snapshot.NameTemplate
	}
	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	name, err := tmpl.New(ctx).Apply(ctx.Config.Snapshot.VersionTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse snapshot name: %w", err)
	}
	if name == "" {
		return fmt.Errorf("empty snapshot name")
	}
	ctx.Version = name
	log.WithField("version", ctx.Version).Infof("building snapshot...")
	return nil
}
