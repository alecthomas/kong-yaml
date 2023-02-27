package kongyaml

import (
	"net/netip"
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
		Dict          map[string]string
		NestedDict    map[string]map[string]bool
		NonStringDict map[netip.Addr][]string
		TypedDict     map[string]struct {
			Foo string
			Bar float64
		}
		TypedSlice []struct{ Foo string }
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
dict:
    foo: bar
nested-dict:
    foo:
        bar: true # also settable as --nested-dict=foo=bar=true
non-string-dict:
    "1.2.3.4": ["foo", "bar"]
typed-dict:
    foo:
        foo: bar
        bar: 1.337
typed-slice:
    - foo: bar
    - foo: baz
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
		Dict: map[string]string{
			"foo": "bar",
		},
		NestedDict: map[string]map[string]bool{
			"foo": {
				"bar": true,
			},
		},
		NonStringDict: map[netip.Addr][]string{
			netip.MustParseAddr("1.2.3.4"): {"foo", "bar"},
		},
		TypedDict: map[string]struct {
			Foo string
			Bar float64
		}{
			"foo": {
				Foo: "bar",
				Bar: 1.337,
			},
		},
		TypedSlice: []struct{ Foo string }{
			{Foo: "bar"},
			{Foo: "baz"},
		},
	}
	require.Equal(t, expected, cli)
}

func TestEmptyFile(t *testing.T) {
	type CLI struct {
		FlagName string
	}
	var cli CLI
	r := strings.NewReader("")
	resolver, err := Loader(r)
	require.NoError(t, err)
	parser, err := kong.New(&cli, kong.Resolvers(resolver))
	require.NoError(t, err)
	_, err = parser.Parse([]string{})
	require.NoError(t, err)
	expected := CLI{
		FlagName: "",
	}
	require.Equal(t, expected, cli)
}
