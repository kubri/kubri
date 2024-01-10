package pipe

import (
	"errors"

	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/azureblob"
	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/gcs"
	"github.com/abemedia/appcast/target/github"
	"github.com/abemedia/appcast/target/s3"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

var errInvalidTarget = errors.New("target: invalid type")

type targetConfig struct {
	*azureblobTarget
	*gcsTarget
	*s3Target
	*fileTarget
	*githubTarget
}

func (tc *targetConfig) UnmarshalYAML(node *yaml.Node) error {
	var typ struct {
		Type string `yaml:"type"`
	}
	if err := node.Decode(&typ); err != nil {
		return err
	}

	switch typ.Type {
	case "azureblob":
		return node.Decode(&tc.azureblobTarget)
	case "gcs":
		return node.Decode(&tc.gcsTarget)
	case "s3":
		return node.Decode(&tc.s3Target)
	case "file":
		return node.Decode(&tc.fileTarget)
	case "github":
		return node.Decode(&tc.githubTarget)
	default:
		return errInvalidTarget
	}
}

var _ yaml.Unmarshaler = (*targetConfig)(nil)

func (tc targetConfig) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			withType(tc.azureblobTarget, "azureblob"),
			withType(tc.gcsTarget, "gcs"),
			withType(tc.s3Target, "s3"),
			withType(tc.fileTarget, "file"),
			withType(tc.githubTarget, "github"),
		},
	}
}

type azureblobTarget struct {
	Bucket string `yaml:"bucket"`
	Folder string `yaml:"folder,omitempty"`
	URL    string `yaml:"url,omitempty"`
}

type gcsTarget struct {
	Bucket string `yaml:"bucket"`
	Folder string `yaml:"folder,omitempty"`
	URL    string `yaml:"url,omitempty"`
}

type s3Target struct {
	Bucket     string `yaml:"bucket"`
	Folder     string `yaml:"folder,omitempty"`
	Endpoint   string `yaml:"endpoint,omitempty"`
	Region     string `yaml:"region,omitempty"`
	DisableSSL bool   `yaml:"disable-ssl,omitempty"`
	URL        string `yaml:"url,omitempty"`
}

type fileTarget struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url,omitempty"`
}

type githubTarget struct {
	Owner  string `yaml:"owner"`
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch,omitempty"`
	Folder string `yaml:"folder,omitempty"`
}

func getTarget(c *targetConfig) (target.Target, error) {
	switch {
	case c.azureblobTarget != nil:
		return azureblob.New(azureblob.Config(*c.azureblobTarget))
	case c.gcsTarget != nil:
		return gcs.New(gcs.Config(*c.gcsTarget))
	case c.s3Target != nil:
		return s3.New(s3.Config(*c.s3Target))
	case c.fileTarget != nil:
		return file.New(file.Config(*c.fileTarget))
	case c.githubTarget != nil:
		return github.New(github.Config(*c.githubTarget))
	default:
		return nil, errInvalidTarget
	}
}
