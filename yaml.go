package kongyaml

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

type yamlResolver struct {
	config map[string]any
}

func (y *yamlResolver) Validate(app *kong.Application) error {
	var path []*kong.Node
	pop := func() { path = path[:len(path)-1] }
	valid := map[string]bool{}
	for key := range y.config {
		valid[key] = true
	}
	err := kong.Visit(app, func(node kong.Visitable, next kong.Next) error {
		switch node := node.(type) {
		case *kong.Application:
			path = append(path, node.Node)
			defer pop()
		case *kong.Node:
			path = append(path, node)
			defer pop()
		case *kong.Flag:
			key := keyForFlag(path[len(path)-1], node)
			if find(y.config, key) != nil {
				delete(valid, strings.Join(key, "-"))
			}
		}
		return next(nil)
	})
	if err != nil {
		return err
	}
	if len(valid) > 0 {
		keys := make([]string, 0, len(valid))
		for key := range valid {
			keys = append(keys, key)
		}
		return fmt.Errorf("extra configuration keys: %s", strings.Join(keys, ", "))
	}
	return nil
}

func (y *yamlResolver) Resolve(context *kong.Context, parent *kong.Path, flag *kong.Flag) (any, error) {
	// Build a string path up to this flag.
	path := keyForFlag(parent.Node(), flag)
	return find(y.config, path), nil
}

// Loader is a Kong configuration loader for YAML.
func Loader(r io.Reader) (kong.Resolver, error) {
	decoder := yaml.NewDecoder(r)
	config := map[string]any{}
	err := decoder.Decode(config)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("YAML config decode error: %w", err)
	}
	return &yamlResolver{config: config}, nil
}

func find(config map[string]any, path []string) any {
	if len(path) == 0 {
		return config
	}
	for i := 0; i < len(path); i++ {
		prefix := strings.Join(path[:i+1], "-")
		if child, ok := config[prefix].(map[string]any); ok {
			return find(child, path[i+1:])
		}
	}
	return config[strings.Join(path, "-")]
}

func keyForFlag(parent *kong.Node, flag *kong.Flag) []string {
	var path []string
	for n := parent; n != nil && n.Type != kong.ApplicationNode; n = n.Parent {
		path = append([]string{n.Name}, path...)
	}
	path = append(path, flag.Name)
	path = strings.Split(strings.Join(path, "-"), "-")
	return path
}
