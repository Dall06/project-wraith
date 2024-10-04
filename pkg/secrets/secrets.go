package secrets

import (
	"fmt"
	"os"
	"strconv"
)

// Secrets represents the structure to hold all the secrets
type Secrets struct {
	Secrets struct {
		Jwt       string
		DbData    string
		Response  string
		Password  string
		Cookies   string
		Internals string
		Logs      string
	}
	Sms struct {
		ResetAsset string
		From       string
		AccountSID string
		AuthToken  string
	}
	Mail struct {
		From     string
		Password string
		Host     string
		Port     string
	}
	Notifiers struct {
		Bot struct {
			Token string
			Chat  string
		}
	}
	Encryption struct {
		Logs bool // Change Logs to bool
	}
}

// Load populates secrets from environment variables
func Load() (*Secrets, error) {
	secrets := &Secrets{}

	// Populate secrets from environment variables
	secretVars := map[string]*string{
		"JWT_SECRET":       &secrets.Secrets.Jwt,
		"DB_DATA":          &secrets.Secrets.DbData,
		"RESPONSE_SECRET":  &secrets.Secrets.Response,
		"PASSWORD_SECRET":  &secrets.Secrets.Password,
		"COOKIES_SECRET":   &secrets.Secrets.Cookies,
		"INTERNALS_SECRET": &secrets.Secrets.Internals,
		"LOGS_SECRET":      &secrets.Secrets.Logs,
		"SMS_RESET_ASSET":  &secrets.Sms.ResetAsset,
		"SMS_FROM":         &secrets.Sms.From,
		"SMS_ACCOUNT_SID":  &secrets.Sms.AccountSID,
		"SMS_AUTH_TOKEN":   &secrets.Sms.AuthToken,
		"MAIL_FROM":        &secrets.Mail.From,
		"MAIL_PASSWORD":    &secrets.Mail.Password,
		"MAIL_HOST":        &secrets.Mail.Host,
		"MAIL_PORT":        &secrets.Mail.Port,
		"BOT_TOKEN":        &secrets.Notifiers.Bot.Token,
		"BOT_CHAT":         &secrets.Notifiers.Bot.Chat,
	}

	// Load string secrets
	for key, value := range secretVars {
		*value = os.Getenv(key) // Dereference the pointer and set the value
		if *value == "" {
			return nil, fmt.Errorf("%s is required but not set", key)
		}
	}

	// Load boolean for encryption logs
	encryptionLogsEnv := os.Getenv("ENCRYPTION_LOGS")
	if encryptionLogsEnv == "" {
		return nil, fmt.Errorf("ENCRYPTION_LOGS is required but not set")
	}

	encryptionLogs, err := strconv.ParseBool(encryptionLogsEnv)
	if err != nil {
		return nil, fmt.Errorf("ENCRYPTION_LOGS must be a valid boolean (true/false): %v", err)
	}
	secrets.Encryption.Logs = encryptionLogs

	return secrets, nil
}
