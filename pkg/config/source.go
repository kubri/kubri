package config

import (
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"

	"github.com/kubri/kubri/source"
	"github.com/kubri/kubri/source/azureblob"
	"github.com/kubri/kubri/source/file"
	"github.com/kubri/kubri/source/gcs"
	"github.com/kubri/kubri/source/github"
	"github.com/kubri/kubri/source/gitlab"
	"github.com/kubri/kubri/source/local"
	"github.com/kubri/kubri/source/s3"
)

type sourceConfig struct {
	*azureblobSource
	*gcsSource
	*s3Source
	*fileSource
	*githubSource
	*gitlabSource
	*localSource
}

func (tc *sourceConfig) UnmarshalYAML(node *yaml.Node) error {
	var typ struct {
		Type string `yaml:"type"`
	}
	if err := node.Decode(&typ); err != nil {
		return err
	}

	switch typ.Type {
	case "azureblob":
		return node.Decode(&tc.azureblobSource)
	case "gcs":
		return node.Decode(&tc.gcsSource)
	case "s3":
		return node.Decode(&tc.s3Source)
	case "file":
		return node.Decode(&tc.fileSource)
	case "github":
		return node.Decode(&tc.githubSource)
	case "gitlab":
		return node.Decode(&tc.gitlabSource)
	case "local":
		return node.Decode(&tc.localSource)
	default:
		return nil
	}
}

func (tc sourceConfig) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			withType(tc.azureblobSource, "azureblob"),
			withType(tc.gcsSource, "gcs"),
			withType(tc.s3Source, "s3"),
			withType(tc.fileSource, "file"),
			withType(tc.githubSource, "github"),
			withType(tc.gitlabSource, "gitlab"),
			withType(tc.localSource, "local"),
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
	switch {
	case c.azureblobSource != nil:
		return azureblob.New(azureblob.Config(*c.azureblobSource))
	case c.gcsSource != nil:
		return gcs.New(gcs.Config(*c.gcsSource))
	case c.s3Source != nil:
		return s3.New(s3.Config(*c.s3Source))
	case c.fileSource != nil:
		return file.New(file.Config(*c.fileSource))
	case c.githubSource != nil:
		return github.New(github.Config(*c.githubSource))
	case c.gitlabSource != nil:
		return gitlab.New(gitlab.Config(*c.gitlabSource))
	case c.localSource != nil:
		return local.New(local.Config(*c.localSource))
	default:
		return nil, &Error{Errors: []string{"source.type must be one of [azureblob gcs s3 file github gitlab local]"}}
	}
}
