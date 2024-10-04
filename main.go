package main

import (
	"project-wraith/pkg/config"
	"project-wraith/pkg/consts"
	"project-wraith/pkg/core"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/tools"
	"project-wraith/pkg/secrets"
)

// @title project-wraith
// @description REST Api for user interaction
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

	sct, err := secrets.Load()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Logger.FolderPath, sct.Encryption.Logs, sct.Secrets.Logs)
	err = log.Initialize()
	if err != nil {
		panic(err)
	}

	err = core.Start(cfg, sct, log)
	if err != nil {
		panic(err)
	}
}
