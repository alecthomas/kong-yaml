package kongyaml

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
)

// Loader is a Kong configuration loader for YAML.
func Loader(r io.Reader) (kong.ResolverFunc, error) {
	decoder := yaml.NewDecoder(r)
	config := map[interface{}]interface{}{}
	err := decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (string, error) {
		// Build a string path up to this flag.
		path := []string{}
		for n := parent.Node(); n != nil && n.Type != kong.ApplicationNode; n = n.Parent {
			path = append([]string{n.Name}, path...)
		}
		path = append(path, flag.Name)
		value := find(config, path)
		switch value := value.(type) {
		case string:
			return value, nil
		case nil:
			return "", nil
		default:
			return fmt.Sprintf("%v", value), nil
		}
	}, nil
}

func find(config map[interface{}]interface{}, path []string) interface{} {
	if len(path) == 1 {
		return config[path[0]]
	}
	child, ok := config[path[0]].(map[interface{}]interface{})
	if !ok {
		return nil
	}
	return find(child, path[1:])
}
