package sendgrid

import (
	"fmt"

	"github.com/reverie/configs"
	"github.com/reverie/types"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// successful response status code for sendgrid
const statusOk = 202

var client = sendgrid.NewSendClient(configs.Project.SendGrid.Key)

// sender identities
var (
	anish = mail.NewEmail("Anish Mukherjee", "anish.mukherjee@ezflo.in")
	goro  = mail.NewEmail("Gaurav Singhal", "gaurav.singhal@ezflo.in")
)

// List of template IDs
const (
	emailConfirmation = "d-05e19211d4524eefaad00e51e90a98e8"
	postActivation    = "d-9243e54d3a094c849a73b93d8ff41e67"
	postCompletion    = "d-796afee6aa124c709a6095c62e13c463"
	passwordReset     = "d-1227bf471f6745659774a50a93296c8b"
)

// send message with error handling
func send(message *mail.SGMailV3) error {
	resp, err := client.Send(message)
	if err != nil {
		return err
	}
	if resp.StatusCode != statusOk {
		return fmt.Errorf("Error encountered while sending email :- %v", resp)
	}
	return nil
}

// SendConfirmationEmail sends a email confirmation message to the user's email address upon successful registration
func SendConfirmationEmail(username, email, token string) error {
	message := mail.NewV3Mail()

	message.SetFrom(anish)
	message.SetTemplateID(emailConfirmation)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(username, email),
	}
	personalization.AddTos(tos...)
	personalization.SetDynamicTemplateData("link", fmt.Sprintf("%s/auth/confirm-email?token=%s", configs.Project.SendGrid.BackendEndpoint, token))

	message.AddPersonalizations(personalization)
	return send(message)
}

// SendPasswordResetEmail sends a email containing the new login credentials
func SendPasswordResetEmail(email, password string) error {
	message := mail.NewV3Mail()

	message.SetFrom(anish)
	message.SetTemplateID(passwordReset)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail("Me", email),
	}
	personalization.AddTos(tos...)
	personalization.SetDynamicTemplateData("email", email)
	personalization.SetDynamicTemplateData("password", password)

	message.AddPersonalizations(personalization)
	return send(message)
}

// SendPostActivationEmail sends a email notification to us denoting the start of a post
// After this we need to co-ordinate the delivery of equipment in an end-to-end fashion
func SendPostActivationEmail(post *types.Post, users []types.M) error {
	message := mail.NewV3Mail()

	message.SetFrom(anish)
	message.SetTemplateID(postActivation)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		anish,
		goro,
	}
	personalization.AddTos(tos...)
	personalization.SetDynamicTemplateData("post", post)
	personalization.SetDynamicTemplateData("users", users)

	message.AddPersonalizations(personalization)
	return send(message)
}

// SendPostCompletionEmail sends an email notification denoting the end of a post
// It also sends the amount and the bank account details to the post owner and CC's us
func SendPostCompletionEmail(recipentEmail, recipentName, postName string, amount float64) error {
	message := mail.NewV3Mail()

	message.SetFrom(anish)
	message.SetTemplateID(postCompletion)

	personalization := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(recipentName, recipentEmail),
	}
	ccs := []*mail.Email{
		anish,
		goro,
	}
	personalization.AddTos(tos...)
	personalization.AddCCs(ccs...)
	personalization.SetDynamicTemplateData("name", postName)
	personalization.SetDynamicTemplateData("amount", fmt.Sprintf("%.2f", amount))

	message.AddPersonalizations(personalization)
	return send(message)
}
