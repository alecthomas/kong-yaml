package kongyaml

import (
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

type TestSample struct {
	Name string `yaml:"name"`
	Game string `yaml:"game"`
}

func TestYAMLFileMapper(t *testing.T) {
	var cli struct {
		Sample TestSample `type:"yamlfile"`
	}
	opt := kong.NamedMapper("yamlfile", YAMLFileMapper)
	parser, err := kong.New(&cli, opt)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--sample", "testdata/sample.yaml"})
	require.NoError(t, err)

	want := TestSample{Name: "Lee Sedol", Game: "Go"}
	require.Equal(t, want, cli.Sample)
}

func TestYAMLFileMapperErr(t *testing.T) {
	var cli struct {
		Sample TestSample `type:"yamlfile"`
	}
	opts := []kong.Option{
		kong.NamedMapper("yamlfile", YAMLFileMapper),
		kong.Exit(func(int) { t.Log("EXIT") }),
	}
	parser, err := kong.New(&cli, opts...)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--sample", "testdata/MISSING_FILE.yaml"})
	require.Error(t, err)

	_, err = parser.Parse([]string{"--sample"})
	require.Error(t, err)
}
