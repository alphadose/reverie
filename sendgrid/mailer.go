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
	postActivation    = "d-9243e54d3a094c849a73b93d8ff41e67"
	postCompletion    = "d-796afee6aa124c709a6095c62e13c463"
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

// SendPostActivationEmail sends a email notification to us denoting the start of a post
// After this we need to co-ordinate the delivery of equipment in an end-to-end fashion
func SendPostActivationEmail(postID string) error {
	message := mail.NewV3Mail()

	message.SetFrom(sender)
	message.SetTemplateID(postActivation)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		sender,
		mail.NewEmail("Gaurav Singhal", "gauravsinghal5998@gmail.com"),
	}
	personalization.AddTos(tos...)
	personalization.SetDynamicTemplateData("id", postID)

	message.AddPersonalizations(personalization)
	_, err := client.Send(message)
	return err
}

// SendPostCompletionEmail sends an email notification denoting the end of a post
// It also sends the amount and the bank account details to the post owner and CC's us
func SendPostCompletionEmail(recipentEmail, recipentName, postName string, amount float64) error {
	message := mail.NewV3Mail()

	message.SetFrom(sender)
	message.SetTemplateID(postCompletion)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(recipentName, recipentEmail),
	}
	ccs := []*mail.Email{
		sender,
		mail.NewEmail("Gaurav Singhal", "gauravsinghal5998@gmail.com"),
	}
	personalization.AddTos(tos...)
	personalization.AddCCs(ccs...)
	personalization.SetDynamicTemplateData("name", postName)
	personalization.SetDynamicTemplateData("amount", amount)

	message.AddPersonalizations(personalization)
	_, err := client.Send(message)
	return err
}
