package kongyaml

import (
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

// Loader is a Kong configuration loader for YAML.
func Loader(r io.Reader) (kong.Resolver, error) {
	var config map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}
	var f kong.ResolverFunc = func(_ *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		// Build a string path up to this flag.
		var path []string
		for n := parent.Node(); n != nil && n.Type != kong.ApplicationNode; n = n.Parent {
			path = append(path, n.Name)
		}
		// Shallow copy parsed config.
		config := config
		// Path is in reverse order.
		for i := len(path) - 1; i >= 0; i-- {
			var ok bool
			if config, ok = config[path[i]].(map[string]interface{}); !ok {
				return nil, nil
			}
		}
		name := strings.ReplaceAll(flag.Name, "-", "_")
		return config[name], nil
	}
	return f, nil
}
