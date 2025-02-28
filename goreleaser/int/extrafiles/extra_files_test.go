package extrafiles

import (
	"testing"

	"github.com/goreleaser/goreleaser/v2/int/testctx"
	"github.com/goreleaser/goreleaser/v2/int/testlib"
	"github.com/goreleaser/goreleaser/v2/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: "./testdata/file{{ .Env.ONE }}.golden"},
	}

	ctx := testctx.NewWithCfg(config.Project{
		Env: []string{"ONE=1"},
	})
	files, err := Find(ctx, globs)
	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Equal(t, "testdata/file1.golden", files["file1.golden"])
}

func TestBadTemplate(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: "./testdata/file{{ .Env.NOPE }}.golden"},
	}

	ctx := testctx.New()
	files, err := Find(ctx, globs)
	require.Empty(t, files)
	testlib.RequireTemplateError(t, err)
}

func TestShouldGetSpecificFile(t *testing.T) {
	globs := []config.ExtraFile{
		{},                        // empty glob, will be ignored
		{Glob: "./testdata/sub3"}, // will get a file1.golden as well, but will be overridden
		{Glob: "./testdata/file1.golden"},
	}

	files, err := Find(testctx.New(), globs)
	require.NoError(t, err)
	require.Len(t, files, 1)

	require.Equal(t, "testdata/file1.golden", files["file1.golden"])
}

func TestFailToGetSpecificFile(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: "./testdata/file453.golden"},
	}

	files, err := Find(testctx.New(), globs)
	require.EqualError(t, err, "globbing failed for pattern ./testdata/file453.golden: matching \"./testdata/file453.golden\": file does not exist")
	require.Empty(t, files)
}

func TestShouldGetFilesWithSuperStar(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: "./**/file?.golden"},
	}

	files, err := Find(testctx.New(), globs)
	require.NoError(t, err)
	require.Len(t, files, 3)
	require.Equal(t, "testdata/file2.golden", files["file2.golden"])
	require.Equal(t, "testdata/sub3/file1.golden", files["file1.golden"])
	require.Equal(t, "testdata/sub/file5.golden", files["file5.golden"])
}

func TestShouldGetAllFilesWithGoldenExtension(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: "./testdata/*.golden"},
	}

	files, err := Find(testctx.New(), globs)
	require.NoError(t, err)
	require.Len(t, files, 2)
	require.Equal(t, "testdata/file1.golden", files["file1.golden"])
	require.Equal(t, "testdata/file2.golden", files["file2.golden"])
}

func TestShouldGetAllFilesInsideTestdata(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: "./testdata/*"},
	}

	files, err := Find(testctx.New(), globs)
	require.NoError(t, err)
	require.Len(t, files, 4)
	require.Equal(t, "testdata/sub3/file1.golden", files["file1.golden"])
	require.Equal(t, "testdata/file2.golden", files["file2.golden"])
	require.Equal(t, "testdata/file3.gold", files["file3.gold"])
	require.Equal(t, "testdata/sub/file5.golden", files["file5.golden"])
}

func TestTargetName(t *testing.T) {
	globs := []config.ExtraFile{
		{
			Glob:         "./testdata/file1.golden",
			NameTemplate: "file1_{{.Tag}}.golden",
		},
	}

	ctx := testctx.New(testctx.WithCurrentTag("v1.0.0"))
	files, err := Find(ctx, globs)
	require.NoError(t, err)
	require.Len(t, files, 1)

	require.Equal(t, "testdata/file1.golden", files["file1_v1.0.0.golden"])
}

func TestTargetInvalidNameTemplate(t *testing.T) {
	globs := []config.ExtraFile{
		{
			Glob:         "./testdata/file1.golden",
			NameTemplate: "file1_{{.Env.HONK}}.golden",
		},
	}

	ctx := testctx.New()
	files, err := Find(ctx, globs)
	require.Empty(t, files)
	testlib.RequireTemplateError(t, err)
}

func TestTargetNameMatchesMultipleFiles(t *testing.T) {
	globs := []config.ExtraFile{
		{
			Glob:         "./testdata/*",
			NameTemplate: "file1.golden",
		},
	}

	ctx := testctx.New()
	files, err := Find(ctx, globs)
	require.Empty(t, files)
	require.EqualError(t, err, `failed to add extra_file: "./testdata/*" -> "file1.golden": glob matches multiple files`)
}

func TestTargetNameNoMatches(t *testing.T) {
	globs := []config.ExtraFile{
		{
			Glob:         "./testdata/file1.silver",
			NameTemplate: "file1_{{.Tag}}.golden",
		},
	}

	ctx := testctx.New()
	files, err := Find(ctx, globs)
	require.Empty(t, files)
	require.EqualError(t, err, `globbing failed for pattern ./testdata/file1.silver: matching "./testdata/file1.silver": file does not exist`)
}

func TestGlobEvalsToEmpty(t *testing.T) {
	globs := []config.ExtraFile{
		{Glob: `{{ printf "" }}`},
	}

	ctx := testctx.New()
	files, err := Find(ctx, globs)
	require.Empty(t, files)
	require.NoError(t, err)
}

func TestTargetNameNoGlob(t *testing.T) {
	globs := []config.ExtraFile{
		{NameTemplate: "file1.golden"},
	}

	ctx := testctx.New()
	files, err := Find(ctx, globs)
	require.Empty(t, files)
	require.NoError(t, err)
}
