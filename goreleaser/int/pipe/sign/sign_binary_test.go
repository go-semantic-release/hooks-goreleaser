package sign

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goreleaser/goreleaser/v2/int/artifact"
	"github.com/goreleaser/goreleaser/v2/int/skips"
	"github.com/goreleaser/goreleaser/v2/int/testctx"
	"github.com/goreleaser/goreleaser/v2/int/testlib"
	"github.com/goreleaser/goreleaser/v2/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestBinarySignDescription(t *testing.T) {
	require.NotEmpty(t, BinaryPipe{}.String())
}

func TestBinarySignDefault(t *testing.T) {
	ctx := testctx.NewWithCfg(config.Project{
		BinarySigns: []config.Sign{{}},
	})
	err := BinaryPipe{}.Default(ctx)
	require.NoError(t, err)
	require.Equal(t, "gpg", ctx.Config.BinarySigns[0].Cmd)
	require.Equal(t, defaultSignatureName, ctx.Config.BinarySigns[0].Signature)
	require.Equal(t, []string{"--output", "$signature", "--detach-sig", "$artifact"}, ctx.Config.BinarySigns[0].Args)
	require.Equal(t, "binary", ctx.Config.BinarySigns[0].Artifacts)
}

func TestBinarySignDisabled(t *testing.T) {
	ctx := testctx.NewWithCfg(config.Project{
		BinarySigns: []config.Sign{
			{Artifacts: "none"},
		},
	})
	err := BinaryPipe{}.Run(ctx)
	require.EqualError(t, err, "artifact signing is disabled")
}

func TestBinarySignInvalidOption(t *testing.T) {
	ctx := testctx.NewWithCfg(config.Project{
		BinarySigns: []config.Sign{
			{Artifacts: "archive"},
		},
	})
	err := BinaryPipe{}.Run(ctx)
	require.EqualError(t, err, "invalid list of artifacts to sign: archive")
}

func TestBinarySkip(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		require.True(t, BinaryPipe{}.Skip(testctx.New()))
	})

	t.Run("skip sign", func(t *testing.T) {
		ctx := testctx.New(testctx.Skip(skips.Sign))
		require.True(t, BinaryPipe{}.Skip(ctx))
	})

	t.Run("dont skip", func(t *testing.T) {
		ctx := testctx.NewWithCfg(config.Project{
			BinarySigns: []config.Sign{
				{},
			},
		})
		require.False(t, BinaryPipe{}.Skip(ctx))
	})
}

func TestBinaryDependencies(t *testing.T) {
	ctx := testctx.NewWithCfg(config.Project{
		BinarySigns: []config.Sign{
			{Cmd: "cosign"},
			{Cmd: "gpg2"},
		},
	})
	require.Equal(t, []string{"cosign", "gpg2"}, BinaryPipe{}.Dependencies(ctx))
}

func TestBinarySign(t *testing.T) {
	testlib.CheckPath(t, "gpg")
	doTest := func(tb testing.TB, sign config.Sign) []*artifact.Artifact {
		tb.Helper()
		tmpdir := tb.TempDir()

		ctx := testctx.NewWithCfg(config.Project{
			BinarySigns: []config.Sign{sign},
		})

		require.NoError(tb, os.WriteFile(filepath.Join(tmpdir, "bin1"), []byte("foo"), 0o644))
		require.NoError(tb, os.WriteFile(filepath.Join(tmpdir, "bin2"), []byte("foo"), 0o644))

		ctx.Artifacts.Add(&artifact.Artifact{
			Name: "bin1",
			Path: filepath.Join(tmpdir, "bin1"),
			Type: artifact.Binary,
			Extra: map[string]interface{}{
				artifact.ExtraID: "foo",
			},
		})
		ctx.Artifacts.Add(&artifact.Artifact{
			Name: "bin2",
			Path: filepath.Join(tmpdir, "bin2"),
			Type: artifact.Binary,
			Extra: map[string]interface{}{
				artifact.ExtraID: "bar",
			},
		})

		pipe := BinaryPipe{}
		require.NoError(tb, pipe.Default(ctx))

		for i := range ctx.Config.BinarySigns {
			ctx.Config.BinarySigns[i].Args = append(
				[]string{"--homedir", keyring},
				ctx.Config.BinarySigns[i].Args...,
			)
		}
		require.NoError(tb, pipe.Run(ctx))
		return ctx.Artifacts.
			Filter(artifact.ByType(artifact.Signature)).
			List()
	}

	t.Run("default", func(t *testing.T) {
		sigs := doTest(t, config.Sign{})
		require.Len(t, sigs, 2)
	})

	t.Run("filter", func(t *testing.T) {
		sigs := doTest(t, config.Sign{
			ID:  "bar",
			IDs: []string{"bar"},
		})
		require.Len(t, sigs, 1)
	})
}
