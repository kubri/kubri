package config

import (
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"

	"github.com/kubri/kubri/target"
	"github.com/kubri/kubri/target/azureblob"
	"github.com/kubri/kubri/target/file"
	"github.com/kubri/kubri/target/gcs"
	"github.com/kubri/kubri/target/github"
	"github.com/kubri/kubri/target/s3"
)

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
		return nil
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
	Bucket string `yaml:"bucket"           validate:"required"`
	Folder string `yaml:"folder,omitempty" validate:"omitempty,dirname"`
	URL    string `yaml:"url,omitempty"    validate:"omitempty,http_url"`
}

type gcsTarget struct {
	Bucket string `yaml:"bucket"           validate:"required"`
	Folder string `yaml:"folder,omitempty" validate:"omitempty,dirname"`
	URL    string `yaml:"url,omitempty"    validate:"omitempty,http_url"`
}

type s3Target struct {
	Bucket     string `yaml:"bucket"                validate:"required"`
	Folder     string `yaml:"folder,omitempty"      validate:"omitempty,dirname"`
	Endpoint   string `yaml:"endpoint,omitempty"    validate:"omitempty,fqdn|http_url"`
	Region     string `yaml:"region,omitempty"`
	DisableSSL bool   `yaml:"disable-ssl,omitempty"`
	URL        string `yaml:"url,omitempty"         validate:"omitempty,http_url"`
}

type fileTarget struct {
	Path string `yaml:"path"          validate:"required"`
	URL  string `yaml:"url,omitempty" validate:"omitempty,http_url"`
}

type githubTarget struct {
	Owner  string `yaml:"owner"            validate:"required"`
	Repo   string `yaml:"repo"             validate:"required"`
	Branch string `yaml:"branch,omitempty"`
	Folder string `yaml:"folder,omitempty" validate:"omitempty,dirname"`
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
		return nil, &Error{Errors: []string{"target.type must be one of [azureblob gcs s3 file github]"}}
	}
}
