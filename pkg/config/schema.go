package config

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func Schema() []byte {
	s := schema.Reflect(config{})
	b, _ := json.MarshalIndent(s, "", "  ") //nolint:errchkjson
	return b
}

func withType(v any, t string) *jsonschema.Schema {
	s := schema.Reflect(v)
	s.Version = ""
	s.Properties.AddPairs(orderedmap.Pair[string, *jsonschema.Schema]{
		Key:   "type",
		Value: &jsonschema.Schema{Type: "string", Const: t},
	})
	_ = s.Properties.MoveToFront("type")
	return s
}

//nolint:gochecknoglobals
var schema = &jsonschema.Reflector{Anonymous: true, DoNotReference: true, FieldNameTag: "yaml"}
