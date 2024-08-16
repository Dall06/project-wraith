package core

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"os"
	"os/signal"
	"project-wraith/src/config"
	"project-wraith/src/consts"
	"project-wraith/src/modules/domain"
	"project-wraith/src/modules/gateway"
	"project-wraith/src/modules/rules"
	"project-wraith/src/pkg/apikey"
	"project-wraith/src/pkg/db"
	"project-wraith/src/pkg/logger"
	"project-wraith/src/pkg/mail"
	"project-wraith/src/pkg/sms"
	"project-wraith/src/pkg/tools"
)

func Server(cfg *config.Config, log logger.Logger) error {
	dbClient := db.NewClient(cfg.Database.Uri, cfg.Database.Name)
	err := dbClient.Open()
	if err != nil {
		log.Error("Failed to open db client", err)
		return err
	}

	resetSmsAsset, err := tools.ReadAsset(cfg.Sms.ResetAsset)
	if err != nil {
		log.Error("Failed to read sms reset asset", err)
		return err
	}

	mailer := mail.NewMail(
		cfg.Mail.From,
		cfg.Mail.Password,
		cfg.Mail.Host,
		cfg.Mail.Port)
	smsResetSender := sms.NewTwilio(
		cfg.Sms.From,
		cfg.Sms.AccountSID,
		cfg.Sms.AuthToken,
		resetSmsAsset,
	)

	userRepo := domain.NewUserRepository(dbClient)

	userRule := rules.NewRule(userRepo, cfg.Server.JWTSecret)
	userCtrl := gateway.NewUserController(
		log,
		userRule,
		cfg.Server.JWTSecret,
		cfg.Server.CookiesExpiration)

	resetRule := rules.NewResetRule(userRepo, cfg.Server.JWTSecret)
	resetCtrl := gateway.NewResetController(
		log,
		resetRule,
		userRule,
		cfg.Server.JWTSecret,
		cfg.Server.CookiesExpiration,
		mailer,
		smsResetSender,
		cfg.Redirects.ResetUrl,
	)

	staticsCtrl := gateway.NewStaticsController(log, consts.AppManifest.Version)
	serverApiKey := apikey.CrateApiKey(cfg.Server.KeyWord)

	engine := html.New("./views", ".html")

	serverName := fmt.Sprintf("%s@%s", cfg.Server.Name, consts.AppManifest.Version)
	serverConfig := fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ServerHeader:  cfg.Server.Header,
		AppName:       serverName,
		Views:         engine,
	}
	server := fiber.New(serverConfig)

	Middleware(server, cfg.Server.BasePath, serverApiKey, cfg.Server.JWTSecret, cfg.Server.CookiesSecret)
	EnRoute(server, cfg.Server.BasePath, userCtrl, resetCtrl, staticsCtrl)

	go func() {
		listenOn := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		err := server.Listen(listenOn)
		if err != nil {
			log.Error("Failed to listen on port", err)
			return
		}
	}()

	log.Info(
		"Running api server in %s:%d, with base path %s",
		cfg.Server.Host,
		cfg.Server.Port,
		cfg.Server.BasePath)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info("Shutting down server...")
	err = server.Shutdown()
	if err != nil {
		log.Error("Failed to shutdown", err)
		return err
	}

	log.Info("Successfully shutdown!")
	return nil
}
