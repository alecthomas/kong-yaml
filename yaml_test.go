package kongyaml

import (
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	var cli struct {
		FlagName string
		Command  struct {
			NestedFlag string
		} `cmd:""`
	}
	r := strings.NewReader(`
flag-name: "hello world"
command:
    nested-flag: "nested flag"
    number: 1.0
    int: 12342345234534
`)
	resolver, err := Loader(r)
	require.NoError(t, err)
	parser, err := kong.New(&cli, kong.Resolver(resolver))
	require.NoError(t, err)
	_, err = parser.Parse([]string{"command"})
	require.NoError(t, err)
	require.Equal(t, "hello world", cli.FlagName)
	require.Equal(t, "nested flag", cli.Command.NestedFlag)
}
