package actions

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOver(client *sendgrid.Client, email, name string) error {
	from := mail.NewEmail("i9 Team", "noreply@i9fit.co")
	to := mail.NewEmail(name, email)

	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { max-width: 600px; margin: auto; background: #ffffff; padding: 20px; border-radius: 8px; }
			h1 { color: #333333; }
			p { color: #666666; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>All Done</h1>
			<p>Your i9 Giga Membership has officially ended. Feel free to restart it or reach out to us.</p>
		</div>
	</body>
	</html>
	`

	textContent := "All Done\nYour i9 Giga Membership has officially ended. Feel free to restart it or reach out to us."

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Ended", to, textContent, htmlContent)
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: %v", response.StatusCode)
	}
	return nil
}
