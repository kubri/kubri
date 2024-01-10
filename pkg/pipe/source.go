package pipe

import (
	"errors"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/azureblob"
	"github.com/abemedia/appcast/source/file"
	"github.com/abemedia/appcast/source/gcs"
	"github.com/abemedia/appcast/source/github"
	"github.com/abemedia/appcast/source/gitlab"
	"github.com/abemedia/appcast/source/local"
	"github.com/abemedia/appcast/source/s3"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

var errInvalidSource = errors.New("source: invalid type")

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
		return errInvalidSource
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
	Bucket string `yaml:"bucket"`
	Folder string `yaml:"folder,omitempty"`
	URL    string `yaml:"url,omitempty"`
}

type gcsSource struct {
	Bucket string `yaml:"bucket"`
	Folder string `yaml:"folder,omitempty"`
	URL    string `yaml:"url,omitempty"`
}

type s3Source struct {
	Bucket     string `yaml:"bucket"`
	Folder     string `yaml:"folder,omitempty"`
	Endpoint   string `yaml:"endpoint,omitempty"`
	Region     string `yaml:"region,omitempty"`
	DisableSSL bool   `yaml:"disable-ssl,omitempty"`
	URL        string `yaml:"url,omitempty"`
}

type fileSource struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url,omitempty"`
}

type githubSource struct {
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
}

type gitlabSource struct {
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
	URL   string `yaml:"url,omitempty"`
}

type localSource struct {
	Path    string `yaml:"path"`
	Version string `yaml:"version"`
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
		return nil, errInvalidSource
	}
}
