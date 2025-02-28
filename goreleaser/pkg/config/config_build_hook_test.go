package config

import (
	"testing"

	"github.com/goreleaser/goreleaser/v2/int/yaml"
	"github.com/stretchr/testify/require"
)

func TestBuildHook_justString(t *testing.T) {
	var actual BuildHookConfig

	err := yaml.UnmarshalStrict([]byte(`pre: ./script.sh`), &actual)
	require.NoError(t, err)
	require.Equal(t, Hook{
		Cmd: "./script.sh",
		Env: nil,
	}, actual.Pre[0])
}

func TestBuildHook_stringCmds(t *testing.T) {
	var actual BuildHookConfig

	err := yaml.UnmarshalStrict([]byte(`pre:
 - ./script.sh
 - second-script.sh
`), &actual)
	require.NoError(t, err)

	require.Equal(t, Hooks{
		{
			Cmd: "./script.sh",
			Env: nil,
		},
		{
			Cmd: "second-script.sh",
			Env: nil,
		},
	}, actual.Pre)
}

func TestBuildHook_complex(t *testing.T) {
	var actual BuildHookConfig

	err := yaml.UnmarshalStrict([]byte(`pre:
 - cmd: ./script.sh
   env:
    - TEST=value
`), &actual)
	require.NoError(t, err)
	require.Equal(t, Hook{
		Cmd: "./script.sh",
		Env: []string{"TEST=value"},
	}, actual.Pre[0])
}
