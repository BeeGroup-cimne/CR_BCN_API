package mailing

import (
	"bytes"
	"go-api/internal/helpers/env"
	"html/template"
	"strconv"
)

func CreatePasswordRecoveryMessage(email string, token string, expires int) (string, error) {
	body := `
		<html>
		<body>
			<h3>Password recovery</h3>
			<p>Click the link below to reset your password:</p>
			<p><a href="{{.Website}}/recover-password?token={{.Token}}&email={{.Email}}" target="_blank">Reset Password</a></p>
			<p>This link will expire in {{.Expires}} hours.</p>
		</body>
		</html>
	`
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
