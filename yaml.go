package kongyaml

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
)

// Loader is a Kong configuration loader for YAML.
func Loader(r io.Reader) (kong.Resolver, error) {
	decoder := yaml.NewDecoder(r)
	config := map[interface{}]interface{}{}
	err := decoder.Decode(config)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("YAML config decode error: %w", err)
	}
	return kong.ResolverFunc(func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		// Build a string path up to this flag.
		path := []string{}
		for n := parent.Node(); n != nil && n.Type != kong.ApplicationNode; n = n.Parent {
			path = append([]string{n.Name}, path...)
		}
		path = append(path, flag.Name)
		path = strings.Split(strings.Join(path, "-"), "-")
		return find(config, path), nil
	}), nil
}

func find(config map[interface{}]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return convertToStringMap(config)
	}
	for i := 0; i < len(path); i++ {
		prefix := strings.Join(path[:i+1], "-")
		if child, ok := config[prefix].(map[interface{}]interface{}); ok {
			return find(child, path[i+1:])
		}
	}
	return config[strings.Join(path, "-")]
}

func convertToStringMap(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k.(string)] = v
	}

	return out
}
