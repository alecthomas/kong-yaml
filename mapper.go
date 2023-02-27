package kongyaml

import (
	"os"
	"reflect"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

// YAMLFileMapper implements kong.MapperValue to decode a YAML file into
// a struct field.
//
//	var cli struct {
//	  Profile Profile `type:"yamlfile"`
//	}
//
//	func main() {
//	  kong.Parse(&cli, kong.NamedMapper("yamlfile", YAMLFileMapper))
//	}
var YAMLFileMapper = kong.MapperFunc(decodeYAMLFile) //nolint: gochecknoglobals

func decodeYAMLFile(ctx *kong.DecodeContext, target reflect.Value) error {
	var fname string
	if err := ctx.Scan.PopValueInto("filename", &fname); err != nil {
		return err
	}
	f, err := os.Open(fname) //nolint:gosec
	if err != nil {
		return err
	}
	defer f.Close() //nolint

	return yaml.NewDecoder(f).Decode(target.Addr().Interface())
}
