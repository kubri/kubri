package config

import (
	"errors"

	"github.com/goccy/go-yaml"
	"github.com/invopop/jsonschema"

	"github.com/kubri/kubri/target"
	"github.com/kubri/kubri/target/azureblob"
	"github.com/kubri/kubri/target/file"
	"github.com/kubri/kubri/target/gcs"
	"github.com/kubri/kubri/target/github"
	"github.com/kubri/kubri/target/s3"
)

type targetConfig struct{ Config any }

func (c *targetConfig) UnmarshalYAML(unmarshal func(any) error) error {
	var t struct {
		Type string `yaml:"type"`
	}
	if err := unmarshal(&t); err != nil {
		if _, ok := errors.AsType[*yaml.UnknownFieldError](err); !ok {
			return err
		}
	}

	switch t.Type {
	case "azureblob":
		return decodeTyped[azureblobTarget](unmarshal, &c.Config)
	case "gcs":
		return decodeTyped[gcsTarget](unmarshal, &c.Config)
	case "s3":
		return decodeTyped[s3Target](unmarshal, &c.Config)
	case "file":
		return decodeTyped[fileTarget](unmarshal, &c.Config)
	case "github":
		return decodeTyped[githubTarget](unmarshal, &c.Config)
	default:
		return nil
	}
}

func (targetConfig) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			withType(azureblobTarget{}, "azureblob"),
			withType(gcsTarget{}, "gcs"),
			withType(s3Target{}, "s3"),
			withType(fileTarget{}, "file"),
			withType(githubTarget{}, "github"),
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
	Bucket   string `yaml:"bucket"             validate:"required"`
	Folder   string `yaml:"folder,omitempty"   validate:"omitempty,dirname"`
	Endpoint string `yaml:"endpoint,omitempty" validate:"omitempty,http_url"`
	Region   string `yaml:"region,omitempty"`
	URL      string `yaml:"url,omitempty"      validate:"omitempty,http_url"`
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
	switch c := c.Config.(type) {
	case *azureblobTarget:
		return azureblob.New(azureblob.Config(*c))
	case *gcsTarget:
		return gcs.New(gcs.Config(*c))
	case *s3Target:
		return s3.New(s3.Config(*c))
	case *fileTarget:
		return file.New(file.Config(*c))
	case *githubTarget:
		return github.New(github.Config(*c))
	default:
		return nil, &Error{Errors: []string{"target.type must be one of [azureblob gcs s3 file github]"}}
	}
}
