package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"path/filepath"
)

type Init struct {
	App struct {
		Level string
	}
	Database struct {
		User struct {
			Uri  string
			Name string
		}
		Manager struct {
			Uri  string
			Name string
		}
		License struct {
			Uri  string
			Name string
		}
	}
	Storage struct {
		Bucket string
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
	Options struct {
		EncryptResponse bool
		EncryptDbData   bool
		EncryptLogs     bool
		UploadLogs      bool
		UseLicense      bool
	}
}

func LoadInit(fileName, extension, folderPath string) (*Init, error) {
	fileNameWithExt := fmt.Sprintf("%s.%s", fileName, extension)
	filePath := filepath.Join(folderPath, fileNameWithExt)

	cfgIni, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ini file: %w", err)
	}

	var initConfig Init

	// Database section
	initConfig.Database.User.Uri = cfgIni.Section("database.user").Key("uri").String()
	initConfig.Database.User.Name = cfgIni.Section("database.user").Key("name").String()
	initConfig.Database.Manager.Uri = cfgIni.Section("database.manager").Key("uri").String()
	initConfig.Database.Manager.Name = cfgIni.Section("database.manager").Key("name").String()
	initConfig.Database.License.Uri = cfgIni.Section("database.license").Key("uri").String()
	initConfig.Database.License.Name = cfgIni.Section("database.license").Key("name").String()

	// SMS section
	initConfig.Sms.ResetAsset = cfgIni.Section("sms").Key("reset_asset").String()
	initConfig.Sms.From = cfgIni.Section("sms").Key("from").String()
	initConfig.Sms.AccountSID = cfgIni.Section("sms").Key("account_sid").String()
	initConfig.Sms.AuthToken = cfgIni.Section("sms").Key("auth_token").String()

	// Mail section
	initConfig.Mail.From = cfgIni.Section("mail").Key("from").String()
	initConfig.Mail.Password = cfgIni.Section("mail").Key("password").String()
	initConfig.Mail.Host = cfgIni.Section("mail").Key("host").String()
	initConfig.Mail.Port = cfgIni.Section("mail").Key("port").String()

	// Options section
	initConfig.Options.EncryptResponse = cfgIni.Section("options").Key("encrypt_response").MustBool()
	initConfig.Options.EncryptDbData = cfgIni.Section("options").Key("encrypt_db_data").MustBool()
	initConfig.Options.EncryptLogs = cfgIni.Section("options").Key("encrypt_logs").MustBool()
	initConfig.Options.UploadLogs = cfgIni.Section("options").Key("upload_logs").MustBool()
	initConfig.Options.UseLicense = cfgIni.Section("options").Key("use_license").MustBool()

	return &initConfig, nil
}
