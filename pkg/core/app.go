package core

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"os"
	"os/signal"
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
	"project-wraith/pkg/modules/tools"
	"project-wraith/pkg/secrets"
)

func Start(cfg *config.Config, sct *secrets.Secrets, log logger.Logger) error {
	dbClient := db.NewClient(cfg.Database.Uri, cfg.Database.Name)
	err := dbClient.Open()
	if err != nil {
		log.Error("failed to open db client", err)
		return err
	}

	if consts.LicenseCheck {
		licensesCollection := dbClient.Collection(consts.LicensesCollection)
		licensesCtx := dbClient.Ctx()

		licensesRepo := lics.NewLicenseRepository(*licensesCollection, licensesCtx)
		err = Activate(licensesRepo, cfg.Server.License)
		if err != nil {
			log.Error("failed to activate license", err)
			return err
		}
	}

	resetSmsAsset, err := tools.ReadAsset(sct.Sms.ResetAsset)
	if err != nil {
		log.Error("Failed to read sms reset asset", err)
		return err
	}

	mailer := mail.NewMail(
		sct.Mail.From,
		sct.Mail.Password,
		sct.Mail.Host,
		sct.Mail.Port)
	smsResetSender := sms.NewTwilio(
		sct.Sms.From,
		sct.Sms.AccountSID,
		sct.Sms.AuthToken,
		resetSmsAsset,
	)

	userCollection := dbClient.Collection(consts.UsersCollection)
	userCtx := dbClient.Ctx()

	userRepo := domain.NewUserRepository(*userCollection, userCtx)

	userRule := rules.NewUserRule(
		userRepo, cfg.Options.EncryptDbData, sct.Secrets.DbData, sct.Secrets.Password)
	userCtrl := gateway.NewUserController(
		log,
		userRule,
		sct.Secrets.Jwt,
		cfg.Options.EncryptResponse,
		sct.Secrets.Response,
		cfg.Server.CookiesMinutesLife)

	authCtrl := gateway.NewAuthController(
		log,
		userRule,
		sct.Secrets.Jwt,
		cfg.Server.CookiesMinutesLife)

	resetRule := rules.NewResetRule(userRepo, sct.Secrets.Jwt)
	resetCtrl := gateway.NewResetController(
		log,
		resetRule,
		userRule,
		sct.Secrets.Jwt,
		cfg.Server.CookiesMinutesLife,
		mailer,
		smsResetSender,
		cfg.Redirects.ResetUrl,
	)

	internalsCollection := dbClient.Collection(consts.InternalsCollection)
	internalsCtx := dbClient.Ctx()

	manticore := guard.NewManticore(
		*internalsCollection,
		internalsCtx,
		sct.Secrets.Internals)

	staticsCtrl := gateway.NewStaticsController(log, consts.AppManifest.Version, cfg.Logger.FolderPath, cfg.Server.BasePath)
	serverApiKey := apikey.CrateApiKey(cfg.Server.KeyWord)

	engine := html.New("./public/views", ".html")

	serverName := fmt.Sprintf("%s@%s", cfg.Server.Name, consts.AppManifest.Version)
	serverConfig := fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ServerHeader:  cfg.Server.Header,
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
		fiberApp, log, paths, serverApiKey, sct.Secrets.Jwt, sct.Secrets.Cookies, manticore)
	EnRoute(fiberApp, paths, userCtrl, authCtrl, resetCtrl, staticsCtrl)

	go func() {
		listenOn := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		err := fiberApp.Listen(listenOn)
		if err != nil {
			log.Error("failed to start server", err)
			os.Exit(1)
		}
	}()

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info("shutting down server...")

	err = dbClient.Close()
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

	log.Info("successfully shutdown!")
	return nil
}
