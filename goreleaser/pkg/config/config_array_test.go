package config

import (
	"testing"

	"github.com/goreleaser/goreleaser/v2/int/yaml"
	"github.com/stretchr/testify/require"
)

type Unmarshaled struct {
	Strings StringArray `yaml:"strings,omitempty"`
	Flags   FlagArray   `yaml:"flags,omitempty"`
}

type yamlUnmarshalTestCase struct {
	yaml     string
	expected Unmarshaled
	err      string
}

var stringArrayTests = []yamlUnmarshalTestCase{
	{
		"",
		Unmarshaled{},
		"",
	},
	{
		"strings: []",
		Unmarshaled{
			Strings: StringArray{},
		},
		"",
	},
	{
		"strings: [one two, three]",
		Unmarshaled{
			Strings: StringArray{"one two", "three"},
		},
		"",
	},
	{
		"strings: one two",
		Unmarshaled{
			Strings: StringArray{"one two"},
		},
		"",
	},
	{
		"strings: {key: val}",
		Unmarshaled{},
		"yaml: unmarshal errors:\n  line 1: cannot unmarshal !!map into string",
	},
}

var flagArrayTests = []yamlUnmarshalTestCase{
	{
		"",
		Unmarshaled{},
		"",
	},
	{
		"flags: []",
		Unmarshaled{
			Flags: FlagArray{},
		},
		"",
	},
	{
		"flags: [one two, three]",
		Unmarshaled{
			Flags: FlagArray{"one two", "three"},
		},
		"",
	},
	{
		"flags: one two",
		Unmarshaled{
			Flags: FlagArray{"one", "two"},
		},
		"",
	},
	{
		"flags: {key: val}",
		Unmarshaled{},
		"yaml: unmarshal errors:\n  line 1: cannot unmarshal !!map into string",
	},
}

func TestStringArray(t *testing.T) {
	for _, testCase := range stringArrayTests {
		var actual Unmarshaled

		err := yaml.UnmarshalStrict([]byte(testCase.yaml), &actual)
		if testCase.err == "" {
			require.NoError(t, err)
			require.Equal(t, testCase.expected, actual)
		} else {
			require.EqualError(t, err, testCase.err)
		}
	}
}

func TestFlagArray(t *testing.T) {
	for _, testCase := range flagArrayTests {
		var actual Unmarshaled

		err := yaml.UnmarshalStrict([]byte(testCase.yaml), &actual)
		if testCase.err == "" {
			require.NoError(t, err)
		} else {
			require.EqualError(t, err, testCase.err)
		}
		require.Equal(t, testCase.expected, actual)
	}
}
