package main

import (
	"project-wraith/pkg/config"
	"project-wraith/pkg/consts"
	"project-wraith/pkg/core"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/tools"
)

// @title project-wraith
// @description Golang REST Api for user interaction
func main() {
	manifestPath := tools.BuildPath(consts.ManifestName, consts.ManifestExt, consts.ManifestPath)
	err := consts.ReadManifest(manifestPath)
	if err != nil {
		panic(err)
	}

	cfg, err := config.Load(consts.FileName, consts.Extension, consts.Path)
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Logger.FolderPath)
	err = log.Initialize()
	if err != nil {
		panic(err)
	}

	err = core.Start(cfg, log)
	if err != nil {
		panic(err)
	}
}
