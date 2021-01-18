package sendgrid

import (
	"fmt"

	"github.com/reverie/configs"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	sender = mail.NewEmail("Anish Mukherjee", "anish.mukherjee1996@gmail.com")
	client = sendgrid.NewSendClient(configs.Project.SendGrid.Key)
)

// List of template IDs
var (
	emailConfirmation = "d-05e19211d4524eefaad00e51e90a98e8"
)

// SendConfirmationEmail sends a email confirmation message to the user's email address upon successful registration
func SendConfirmationEmail(username, email, token string) error {
	message := mail.NewV3Mail()

	message.SetFrom(sender)
	message.SetTemplateID(emailConfirmation)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(username, email),
	}
	personalization.AddTos(tos...)
	personalization.SetDynamicTemplateData("link", fmt.Sprintf("%s/auth/confirm-email?token=%s", configs.Project.SendGrid.BackendEndpoint, token))

	message.AddPersonalizations(personalization)
	_, err := client.Send(message)
	return err
}
