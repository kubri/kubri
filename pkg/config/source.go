package config

import (
	"errors"

	"github.com/goccy/go-yaml"
	"github.com/invopop/jsonschema"

	"github.com/kubri/kubri/source"
	"github.com/kubri/kubri/source/azureblob"
	"github.com/kubri/kubri/source/file"
	"github.com/kubri/kubri/source/gcs"
	"github.com/kubri/kubri/source/github"
	"github.com/kubri/kubri/source/gitlab"
	"github.com/kubri/kubri/source/local"
	"github.com/kubri/kubri/source/s3"
)

type sourceConfig struct{ Config any }

func (c *sourceConfig) UnmarshalYAML(unmarshal func(any) error) error {
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
		return decodeTyped[azureblobSource](unmarshal, &c.Config)
	case "gcs":
		return decodeTyped[gcsSource](unmarshal, &c.Config)
	case "s3":
		return decodeTyped[s3Source](unmarshal, &c.Config)
	case "file":
		return decodeTyped[fileSource](unmarshal, &c.Config)
	case "github":
		return decodeTyped[githubSource](unmarshal, &c.Config)
	case "gitlab":
		return decodeTyped[gitlabSource](unmarshal, &c.Config)
	case "local":
		return decodeTyped[localSource](unmarshal, &c.Config)
	default:
		return nil
	}
}

func (sourceConfig) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			withType(azureblobSource{}, "azureblob"),
			withType(gcsSource{}, "gcs"),
			withType(s3Source{}, "s3"),
			withType(fileSource{}, "file"),
			withType(githubSource{}, "github"),
			withType(gitlabSource{}, "gitlab"),
			withType(localSource{}, "local"),
		},
	}
}

type azureblobSource struct {
	Bucket string `yaml:"bucket"           validate:"required"`
	Folder string `yaml:"folder,omitempty" validate:"omitempty,dirname"`
	URL    string `yaml:"url,omitempty"    validate:"omitempty,http_url"`
}

type gcsSource struct {
	Bucket string `yaml:"bucket"           validate:"required"`
	Folder string `yaml:"folder,omitempty" validate:"omitempty,dirname"`
	URL    string `yaml:"url,omitempty"    validate:"omitempty,http_url"`
}

type s3Source struct {
	Bucket   string `yaml:"bucket"             validate:"required"`
	Folder   string `yaml:"folder,omitempty"   validate:"omitempty,dirname"`
	Endpoint string `yaml:"endpoint,omitempty" validate:"omitempty,http_url"`
	Region   string `yaml:"region,omitempty"`
	URL      string `yaml:"url,omitempty"      validate:"omitempty,http_url"`
}

type fileSource struct {
	Path string `yaml:"path"          validate:"required,dir"`
	URL  string `yaml:"url,omitempty" validate:"omitempty,http_url"`
}

type githubSource struct {
	Owner string `yaml:"owner" validate:"required"`
	Repo  string `yaml:"repo"  validate:"required"`
}

type gitlabSource struct {
	Owner string `yaml:"owner"         validate:"required"`
	Repo  string `yaml:"repo"          validate:"required"`
	URL   string `yaml:"url,omitempty" validate:"omitempty,http_url"`
}

type localSource struct {
	Path    string `yaml:"path"    validate:"required,dir"`
	Version string `yaml:"version" validate:"required,version"`
}

func getSource(c *sourceConfig) (*source.Source, error) {
	switch c := c.Config.(type) {
	case *azureblobSource:
		return azureblob.New(azureblob.Config(*c))
	case *gcsSource:
		return gcs.New(gcs.Config(*c))
	case *s3Source:
		return s3.New(s3.Config(*c))
	case *fileSource:
		return file.New(file.Config(*c))
	case *githubSource:
		return github.New(github.Config(*c))
	case *gitlabSource:
		return gitlab.New(gitlab.Config(*c))
	case *localSource:
		return local.New(local.Config(*c))
	default:
		return nil, &Error{Errors: []string{"source.type must be one of [azureblob gcs s3 file github gitlab local]"}}
	}
}
