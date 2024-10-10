# Project Wraith

This is a Go project built using the Fiber framework, a fast and lightweight web framework for Go. This project demonstrates basic usage of Fiber for building RESTful APIs.

## Features

- User authentication and session management
- CRUD operations for user management
- Password reset functionality
- JSON and HTML responses
- Integration with external services (e.g., email and SMS)

## Installation

1. Clone the repository:

   git clone https://github.com/yourusername/your-repo.git
   cd your-repo

2. Install dependencies:

   Ensure you have Go installed. Then, run:

   go mod tidy

3. Run the application:

   go run main.go

   The application will start on localhost:8080 by default.

## Run Swagger

1. Run the Swagger CLI:

   go install github.com/swaggo/swag/cmd/swag@latest
   swag init

## Configuration Files

### `config.env`

# Server configuration environment variables
SERVER_KEY_WORD = your_keyword

# Server secrets environment variables
SECRET_JWT = your_jwt_secret
SECRET_DB = your_db_secret
SECRET_RESPONSE = your_response_secret
SECRET_PASSWORD = your_password_secret
SECRET_COOKIES = your_cookies_secret
SECRET_INTERNALS = your_internals_secret
SECRET_LOGS = your_logs_secret

# Server notifiers environment variables
NOTIFIER_TLG_BOT_TOKEN = your_bot_token
NOTIFIER_TLG_BOT_CHAT = your_bot_chat

### `config.ini`

[database.user]
uri = db_uri
name = dbname

[database.manager]
uri = db_uri
name = dbname-manager

[database.license]
uri = db_uri
name = dbname-license

[sms]
resetAsset = ./public/texts/reset_sms.txt
from = 477XXXXXXXXX
accountSID = your_account_sid
authToken = your_auth_token

[mail]
from = your_email_from
password = your_email_password
host = smtp.gmail.com
port = 587

[options]
encrypt_response = true
encrypt_db_data = true
encrypt_logs = true
upload_logs = true
use_license = true

### `config.lic`

im_your_license

### `config.yaml`

server:
host: "localhost"
port: 8080
env: "development"
basePath: "/project-wraith/api/v1"
cookiesMinutesLife: 15

logger:
debug: true
folderPath: "./logs"

redirects:
resetUrl: "http://localhost:8080/reset"
