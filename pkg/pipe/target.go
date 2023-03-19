package pipe

import (
	"fmt"

	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/azureblob"
	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/gcs"
	"github.com/abemedia/appcast/target/github"
	"github.com/abemedia/appcast/target/s3"
	"github.com/mitchellh/mapstructure"
)

type targetConfig map[string]any

func getTarget(c targetConfig) (target.Target, error) {
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
	default:
		return nil, fmt.Errorf("invalid target type: %s", c["type"])
	}
}
