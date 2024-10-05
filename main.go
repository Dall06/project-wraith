package main

import (
	"project-wraith/pkg/config"
	"project-wraith/pkg/consts"
	"project-wraith/pkg/core"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/tools"
)

// @title project-wraith
// @description REST Api for user interaction
func main() {
	manifestPath := tools.BuildPath(
		consts.ManifestName,
		consts.ManifestExt,
		consts.ManifestPath)
	err := consts.ReadManifest(manifestPath)
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadSetup(
		consts.SetupFileName,
		consts.SetupExtension,
		consts.SetupPath)
	if err != nil {
		panic(err)
	}

	ini, err := config.LoadInit(
		consts.InitFileName,
		consts.InitExtension,
		consts.InitPath)
	if err != nil {
		panic(err)
	}

	sct, err := config.LoadSecrets(
		consts.SecretsFileName,
		consts.SecretsExtension,
		consts.SecretsPath)
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Logger.FolderPath, ini.Options.EncryptLogs, sct.Secrets.Logs)
	err = log.Initialize()
	if err != nil {
		panic(err)
	}

	err = core.Start(cfg, sct, ini, log)
	if err != nil {
		panic(err)
	}
}
