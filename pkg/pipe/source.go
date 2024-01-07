package pipe

import (
	"fmt"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/azureblob"
	"github.com/abemedia/appcast/source/file"
	"github.com/abemedia/appcast/source/gcs"
	"github.com/abemedia/appcast/source/github"
	"github.com/abemedia/appcast/source/gitlab"
	"github.com/abemedia/appcast/source/local"
	"github.com/abemedia/appcast/source/s3"
	"github.com/mitchellh/mapstructure"
)

type sourceConfig map[string]any

func getSource(c sourceConfig) (*source.Source, error) {
	switch c["type"] {
	case "azureblob":
		opt := &azureblob.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return azureblob.New(*opt)
	case "gcs":
		opt := &gcs.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return gcs.New(*opt)
	case "s3":
		opt := &s3.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return s3.New(*opt)
	case "file":
		opt := &file.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return file.New(*opt)
	case "github":
		opt := &github.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return github.New(*opt)
	case "gitlab":
		opt := &gitlab.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return gitlab.New(*opt)
	case "local":
		opt := &local.Config{}
		if err := mapstructure.Decode(c, opt); err != nil {
			return nil, err
		}
		return local.New(*opt)
	default:
		return nil, fmt.Errorf("invalid source type: %s", c["type"])
	}
}
