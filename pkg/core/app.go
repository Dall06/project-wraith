package core

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"project-wraith/pkg/config"
	"project-wraith/pkg/consts"
	"project-wraith/pkg/internal/domain"
	"project-wraith/pkg/internal/gateway"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/apikey"
	"project-wraith/pkg/modules/db"
	"project-wraith/pkg/modules/guard"
	"project-wraith/pkg/modules/lics"
	"project-wraith/pkg/modules/link"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/mail"
	"project-wraith/pkg/modules/sms"
	"project-wraith/pkg/modules/storage"
	"project-wraith/pkg/modules/tools"
)

func Start(cfg *config.Setup, sct *config.Secrets, ini *config.Init, log logger.Logger) error {
	userDbClient := db.NewClient(ini.Database.User.Uri, ini.Database.User.Name)
	err := userDbClient.Open()
	if err != nil {
		log.Error("failed to open db client", err)
		return err
	}

	licenseDbClient := db.NewClient(ini.Database.License.Uri, ini.Database.License.Name)
	err = licenseDbClient.Open()
	if err != nil {
		log.Error("failed to open db client", err)
		return err
	}

	managerDbClient := db.NewClient(ini.Database.Manager.Uri, ini.Database.Manager.Name)
	err = managerDbClient.Open()
	if err != nil {
		log.Error("failed to open db client", err)
		return err
	}

	if ini.Options.UseLicense {
		licString, err := config.LoadLicense(consts.LicFileName, consts.LicExtension, consts.LicPath)
		if err != nil {
			log.Error("failed to load license", err)
			return err
		}

		licensesCollection := licenseDbClient.Collection(consts.LicensesCollection)
		licensesCtx := licenseDbClient.Ctx()

		licensesRepo := lics.NewLicenseRepository(*licensesCollection, licensesCtx)
		err = Activate(licensesRepo, licString)
		if err != nil {
			log.Error("failed to activate license", err)
			return err
		}
	}

	resetSmsAsset, err := tools.ReadAsset(ini.Sms.ResetAsset)
	if err != nil {
		log.Error("Failed to read sms reset asset", err)
		return err
	}

	mailer := mail.NewMail(
		ini.Mail.From,
		ini.Mail.Password,
		ini.Mail.Host,
		ini.Mail.Port)
	smsResetSender := sms.NewTwilio(
		ini.Sms.From,
		ini.Sms.AccountSID,
		ini.Sms.AuthToken,
		resetSmsAsset,
	)

	userCollection := userDbClient.Collection(consts.UsersCollection)
	userCtx := userDbClient.Ctx()

	userRepo := domain.NewUserRepository(*userCollection, userCtx)

	userRule := rules.NewUserRule(
		userRepo, ini.Options.EncryptDbData, sct.Keys.DbData, sct.Keys.Password)
	userCtrl := gateway.NewUserController(
		log,
		userRule,
		sct.Keys.Jwt,
		ini.Options.EncryptResponse,
		sct.Keys.Response,
		cfg.Server.CookiesMinutesLife)

	authCtrl := gateway.NewAuthController(
		log,
		userRule,
		sct.Keys.Jwt,
		cfg.Server.CookiesMinutesLife)

	resetRule := rules.NewResetRule(userRepo, sct.Keys.Jwt)
	resetCtrl := gateway.NewResetController(
		log,
		resetRule,
		userRule,
		sct.Keys.Jwt,
		cfg.Server.CookiesMinutesLife,
		mailer,
		smsResetSender,
		cfg.Redirects.ResetUrl,
	)

	internalsCollection := managerDbClient.Collection(consts.InternalsCollection)
	internalsCtx := managerDbClient.Ctx()

	manticore := guard.NewManticore(
		*internalsCollection,
		internalsCtx,
		sct.Keys.Internals)

	staticsCtrl := gateway.NewStaticsController(log, consts.AppManifest.Version, cfg.Logger.FolderPath, cfg.Server.BasePath)
	serverApiKey := apikey.CrateApiKey(sct.Server.KeyWord)

	engine := html.New("./public/views", ".html")

	serverName := fmt.Sprintf("%s@%s", consts.ServerName, consts.AppManifest.Version)
	serverConfig := fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ServerHeader:  consts.ServerHeader,
		AppName:       serverName,
		Views:         engine,
		ErrorHandler:  link.Error,
	}

	fiberApp := fiber.New(serverConfig)
	paths := map[string]string{
		"user":    fmt.Sprintf("%s/user", cfg.Server.BasePath),
		"auth":    fmt.Sprintf("%s/auth", cfg.Server.BasePath),
		"reset":   fmt.Sprintf("%s/reset", cfg.Server.BasePath),
		"hello":   fmt.Sprintf("%s/hello", cfg.Server.BasePath),
		"swagger": fmt.Sprintf("%s/swagger/*", cfg.Server.BasePath),
		"logs":    fmt.Sprintf("%s/logs", cfg.Server.BasePath),
		"metrics": fmt.Sprintf("%s/metrics", cfg.Server.BasePath),
	}

	Middleware(
		fiberApp, log, paths, serverApiKey, sct.Keys.Jwt, sct.Keys.Cookies, manticore)
	EnRoute(fiberApp, paths, userCtrl, authCtrl, resetCtrl, staticsCtrl)

	listenOn := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	err = fiberApp.Listen(listenOn)
	if err != nil {
		log.Error("failed to start server", err)
		return err
	}

	startLog := fmt.Sprintf(
		"Running API server on %s:%d with base path %s",
		cfg.Server.Host, cfg.Server.Port, cfg.Server.BasePath)
	log.Info(startLog)

	fmt.Printf(
		"\nHelpful routes:\n"+
			"%s:%d%s/hello\n"+
			"%s:%d%s/metrics\n"+
			"%s:%d%s/log\n",
		cfg.Server.Host, cfg.Server.Port, cfg.Server.BasePath, // For /hello
		cfg.Server.Host, cfg.Server.Port, cfg.Server.BasePath, // For /metrics
		cfg.Server.Host, cfg.Server.Port, cfg.Server.BasePath, // For /log
	)

	log.Info("shutting down server...")

	err = userDbClient.Close()
	if err != nil {
		log.Error("failed to close db client", err)
		return err
	}
	err = managerDbClient.Close()
	if err != nil {
		log.Error("failed to close db client", err)
		return err
	}
	err = licenseDbClient.Close()
	if err != nil {
		log.Error("failed to close db client", err)
		return err
	}

	log.Info("successfully closed db client!")

	err = fiberApp.Shutdown()
	if err != nil {
		log.Error("failed to shutdown", err)
		return err
	}

	log.Info("successfully shutdown server!")

	return nil
}

func Teardown(cfg *config.Setup, sct *config.Secrets, ini *config.Init) error {
	objectStorage := storage.NewObjectStorage(
		sct.Storage.AccessKey,
		sct.Storage.SecretKey,
	)

	type uploadInfo struct {
		bucket     string
		directory  string
		filename   string
		permission string
		encrypt    bool
		encryptKey string
	}

	var toUpload []uploadInfo

	if ini.Options.UploadLogs {
		logFiles := []string{"info.log", "warn.log", "error.log"}

		for _, logFile := range logFiles {
			directory := ""
			if logFile == "info.log" {
				directory = fmt.Sprintf("%s/%s", cfg.Logger.FolderPath)
			}

			toUpload = append(toUpload, uploadInfo{
				bucket:     ini.Storage.Bucket,
				directory:  directory,
				filename:   fmt.Sprintf("%s/%s", cfg.Logger.FolderPath, logFile),
				permission: "public-read",
				encrypt:    ini.Options.EncryptLogs,
				encryptKey: sct.Keys.Logs,
			})
		}
	}

	for _, upload := range toUpload {
		_, err := storage.UploadObject(
			objectStorage,
			upload.bucket,
			upload.directory,
			upload.filename,
			upload.permission,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
