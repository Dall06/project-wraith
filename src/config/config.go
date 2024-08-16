package config

type Config struct {
	Server struct {
		Host              string
		Port              int
		Name              string
		Header            string
		Env               string
		BasePath          string
		Debug             bool
		JWTSecret         string
		KeyWord           string
		CookiesSecret     string
		CookiesExpiration int
	}
	Database struct {
		Uri  string
		Name string
	}
	Logger struct {
		FolderPath string
	}
	Encryption struct {
		Key string
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
	Redirects struct {
		ResetUrl string
	}
}
