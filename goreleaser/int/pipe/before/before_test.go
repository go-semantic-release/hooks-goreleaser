package before

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caarlos0/log"
	"github.com/goreleaser/goreleaser/v2/int/skips"
	"github.com/goreleaser/goreleaser/v2/int/testctx"
	"github.com/goreleaser/goreleaser/v2/int/testlib"
	"github.com/goreleaser/goreleaser/v2/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	code := m.Run()
	log.SetLevel(log.InfoLevel)
	os.Exit(code)
}

func TestDescription(t *testing.T) {
	require.NotEmpty(t, Pipe{}.String())
}

func TestRunPipe(t *testing.T) {
	for _, tc := range [][]string{
		nil,
		{},
		{"go version"},
		{"go version", "go list"},
		{`bash -c "go version; echo \"lala spaces and such\""`},
	} {
		ctx := testctx.NewWithCfg(
			config.Project{
				Before: config.Before{
					Hooks: tc,
				},
			},
		)
		require.NoError(t, Pipe{}.Run(ctx))
	}
}

func TestRunPipeInvalidCommand(t *testing.T) {
	ctx := testctx.NewWithCfg(
		config.Project{
			Before: config.Before{
				Hooks: []string{`bash -c "echo \"unterminated command\"`},
			},
		},
	)
	require.EqualError(t, Pipe{}.Run(ctx), "invalid command line string")
}

func TestRunPipeFail(t *testing.T) {
	for _, tc := range []string{
		"go tool foobar",
		"sh ./testdata/foo.sh",
	} {
		ctx := testctx.NewWithCfg(
			config.Project{
				Before: config.Before{
					Hooks: []string{tc},
				},
			},
		)
		err := Pipe{}.Run(ctx)
		require.ErrorContains(t, err, "hook failed")
	}
}

func TestRunWithEnv(t *testing.T) {
	f := filepath.Join(t.TempDir(), "testfile")
	require.NoError(t, Pipe{}.Run(testctx.NewWithCfg(
		config.Project{
			Env: []string{
				"TEST_FILE=" + f,
			},
			Before: config.Before{
				Hooks: []string{"touch {{ .Env.TEST_FILE }}"},
			},
		},
	)))
	require.FileExists(t, f)
}

func TestInvalidTemplate(t *testing.T) {
	testlib.RequireTemplateError(t, Pipe{}.Run(testctx.NewWithCfg(
		config.Project{
			Before: config.Before{
				Hooks: []string{"touch {{ .fasdsd }"},
			},
		},
	)))
}

func TestSkip(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		require.True(t, Pipe{}.Skip(testctx.New()))
	})

	t.Run("skip before", func(t *testing.T) {
		ctx := testctx.NewWithCfg(config.Project{
			Before: config.Before{
				Hooks: []string{""},
			},
		}, testctx.Skip(skips.Before))
		require.True(t, Pipe{}.Skip(ctx))
	})

	t.Run("dont skip", func(t *testing.T) {
		ctx := testctx.NewWithCfg(config.Project{
			Before: config.Before{
				Hooks: []string{""},
			},
		})
		require.False(t, Pipe{}.Skip(ctx))
	})
}
