package mailing

import (
	"bytes"
	"go-api/internal/helpers/env"
	"html/template"
	"strconv"
)

func CreateAccessRequestAcceptedMessage(email string, token string, expires int) (string, error) {
	body := `<html>
		<body>
			<h3>Welcome to ClimateReady Barcelona</h3>
			<p>Your request to access the platform has been granted. Visit us on <a href="{{.Website}}" target="_blank">{{.Website}}</a>, but first...</p>
			<p>Click the link below to set your password:</p>
			<p><a href="{{.Website}}/recover-password?token={{.Token}}&email={{.Email}}&newuser=true" target="_blank">Set new password.</a></p>
			<p>This link will expire in {{.Expires}} hours.</p>
		</body>
		</html>`

	data := struct {
		Token   string
		Email   string
		Expires string
		Website string
	}{
		Token:   token,
		Email:   email,
		Expires: strconv.Itoa(expires),
		Website: env.GetVar("MAILING_WEBSITE"),
	}

	var buf bytes.Buffer
	tmpl := template.Must(template.New("email").Parse(body))
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
