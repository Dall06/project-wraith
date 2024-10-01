package consts

import (
	"github.com/pkg/errors"
	"project-wraith/pkg/modules/tools"
)

var AppManifest *Manifest

const (
	ManifestName = "manifest"
	ManifestExt  = "yaml"
	ManifestPath = "."
)

type Manifest struct {
	Version     string `yaml:"version"`
	Name        string `yaml:"name"`
	Env         string `yaml:"env"`
	Description string `yaml:"description"`
	ReleaseDate string `yaml:"releaseDate"`
	Author      string `yaml:"author"`
}

func ReadManifest(filePath string) error {
	err := tools.ReadYaml(filePath, &AppManifest)
	if err != nil {
		return errors.Errorf("failed to read manifest: %v", err)
	}

	return nil
}
