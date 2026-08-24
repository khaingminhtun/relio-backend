package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
)

func renderTemplate(
	job redis.EmailJob,
) (string, error) {

	templatePath := fmt.Sprintf(
		"internal/shared/email/templates/%s.html",
		job.Template,
	)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf(
			"parse email template: %w",
			err,
		)
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, job.Data); err != nil {
		return "", fmt.Errorf(
			"execute email template: %w",
			err,
		)
	}

	return body.String(), nil
}
