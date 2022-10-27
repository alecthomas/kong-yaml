# Kong YAML utilities [![](https://godoc.org/github.com/alecthomas/kong-yaml?status.svg)](http://godoc.org/github.com/alecthomas/kong-yaml) [![CircleCI](https://img.shields.io/circleci/project/github/alecthomas/kong-yaml.svg)](https://circleci.com/gh/alecthomas/kong-yaml)

## Configuration loader

Use it like so:

```go
parser, err := kong.New(&cli, kong.Configuration(kongyaml.Loader, "/etc/myapp/config.yaml", "~/.myapp.yaml"))
```

## YAMLFileMapper

YAMLFileMapper implements kong.MapperValue to decode a YAML file into
a struct field.

Use it like so:

```go
var cli struct {
  Profile Profile `type:"yamlfile"`
}

func main() {
  kong.Parse(&cli, kong.NamedMapper("yamlfile", kongyaml.YAMLFileMapper))
}
``` 
