package config

type Config struct {
	Server struct {
		Host               string
		Port               int
		Name               string
		Header             string
		Env                string
		BasePath           string
		KeyWord            string
		CookiesMinutesLife int
		License            string
	}
	Database struct {
		Uri  string
		Name string
	}
	Logger struct {
		Debug      bool
		FolderPath string
	}
	Options struct {
		EncryptResponse bool
		EncryptDbData   bool
	}
	Secrets struct {
		Jwt       string
		DbData    string
		Response  string
		Password  string
		Cookies   string
		Internals string
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
	Notifiers struct {
		Bot struct {
			Token string
			Chat  string
		}
	}
}
