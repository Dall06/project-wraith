package sms

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"project-wraith/src/pkg/req"
	"project-wraith/src/pkg/tools"
	"strings"
)

type Twilio interface {
	SendSMSTwilio(to string, useAsset bool, args ...string) (string, error)
}

type twilio struct {
	from       string
	accountSID string
	authToken  string
	asset      string
}

func NewTwilio(from, accountSID, authToken string, asset string) Twilio {
	return &twilio{
		from:       from,
		accountSID: accountSID,
		authToken:  authToken,
		asset:      asset,
	}
}

func (tc twilio) SendSMSTwilio(to string, useAsset bool, args ...string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", tc.accountSID)

	var message string
	if useAsset {
		message = tools.FormatAssetContent(tc.asset, args)
	} else {
		message = strings.Join(args, " ")
	}

	data := url.Values{}
	data.Set("From", tc.from)
	data.Set("To", to)
	data.Set("Body", message)

	authStr := tc.accountSID + ":" + tc.authToken
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))

	smsRequest := req.HTTPRequest{
		Method: "POST",
		URL:    apiURL,
		Headers: map[string]string{
			"Authorization": "Basic " + encodedAuth,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Body: []byte(data.Encode()),
	}

	return req.SendRequest(smsRequest)
}
