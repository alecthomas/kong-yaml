package kongyaml

import (
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	type CLI struct {
		FlagName string
		Names    []string
		Command  struct {
			NestedFlag string
		} `cmd:""`
		Embedded struct {
			One string
			Two bool
		} `embed:"" prefix:"embed-"`
	}
	var cli CLI
	r := strings.NewReader(`
flag-name: "hello world"
embed:
    one: "str"
    two: true
names:
    - "one"
    - "two"
    - "three"
command:
    nested-flag: "nested flag"
    number: 1.0
    int: 12342345234534
`)
	resolver, err := Loader(r)
	require.NoError(t, err)
	parser, err := kong.New(&cli, kong.Resolvers(resolver))
	require.NoError(t, err)
	_, err = parser.Parse([]string{"command"})
	require.NoError(t, err)
	expected := CLI{
		FlagName: "hello world",
		Names:    []string{"one", "two", "three"},
		Command: struct {
			NestedFlag string
		}{NestedFlag: "nested flag"},
		Embedded: struct {
			One string
			Two bool
		}{
			One: "str",
			Two: true,
		},
	}
	require.Equal(t, expected, cli)
}
