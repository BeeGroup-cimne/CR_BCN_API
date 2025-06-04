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
			<h3>Welcome to Climate Ready Barcelona</h3>
			<p>You have been granted access to the website. Visit us on <a href="{{.Website}}" target="_blank">{{.Website}}</a>, but first...</p>
			<p>These are your new login credentials:<br />
			<strong>Username</strong>: {{.Email}}<br />
			<strong>Password</strong>: <a href="{{.Website}}/recover-password?token={{.Token}}&email={{.Email}}&newuser=true" target="_blank">Click this link to set the password</a><br />
			This link will expire in {{.Expires}} hours.</p>
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
