package hooks

import (
	"os"

	"github.com/goreleaser/goreleaser/v2/pkg/config"
)

var goReleaserConfigPaths = [6]string{
	".config/goreleaser.yml",
	".config/goreleaser.yaml",
	".goreleaser.yml",
	".goreleaser.yaml",
	"goreleaser.yml",
	"goreleaser.yaml",
}

func readGoReleaserConfig() (config.Project, error) {
	for _, file := range goReleaserConfigPaths {
		projectConfig, err := config.Load(file)
		if err == nil {
			return projectConfig, nil
		}
		if os.IsNotExist(err) {
			continue
		}
		return projectConfig, err
	}
	return config.Project{}, nil
}
