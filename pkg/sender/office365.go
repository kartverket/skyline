package sender

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/kartverket/skyline/pkg/email"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/pkg/errors"
	"log/slog"
	"net/mail"
)

type office365sender struct {
	graphClient *msgraphsdk.GraphServiceClient
}

func NewOffice365Sender(tenantId string, clientId string, clientSecret string) (Sender, error) {
	slog.Info("Creating new client", "tenant-id", tenantId, "client-id", clientId)

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, err
	}

	c, err := msgraphsdk.NewGraphServiceClientWithCredentials(credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}

	return &office365sender{c}, nil
}

func (s *office365sender) Send(ctx context.Context, email *email.SkylineEmail) error {
	payload, err := mapToGraphMail(email)
	if err != nil {
		return errors.Wrap(err, "could not convert received message to fit within Microsoft Graph API models")
	}

	err = s.graphClient.Me().SendMail().Post(ctx, payload, nil)
	return err
}

func mapToGraphMail(email *email.SkylineEmail) (*graphusers.ItemSendMailPostRequestBody, error) {
	// Outer payload
	requestBody := graphusers.NewItemSendMailPostRequestBody()
	message := graphmodels.NewMessage()

	// Subject
	message.SetSubject(&email.Headers.Subject)

	// From
	message.SetFrom(toGraphRecipient(email.Headers.From...)[0])

	// Recipients
	message.SetToRecipients(toGraphRecipient(email.Headers.To...))
	message.SetCcRecipients(toGraphRecipient(email.Headers.Cc...))
	message.SetBccRecipients(toGraphRecipient(email.Headers.Bcc...))

	// Metadata
	message.SetSentDateTime(&email.Headers.Date)

	// TODO: Attachments
	// TODO: IDs

	// Body
	body := graphmodels.NewItemBody()
	if email.IsPlaintext() {
		contentType := graphmodels.TEXT_BODYTYPE
		body.SetContentType(&contentType)
		body.SetContent(&email.Text)
	} else if email.IsHTML() {
		contentType := graphmodels.HTML_BODYTYPE
		body.SetContentType(&contentType)
		body.SetContent(&email.HTML)
	} else {
		return nil, errors.New("unsupported content type: " + email.Headers.ContentType.ContentType)
	}
	message.SetBody(body)

	requestBody.SetMessage(message)
	saveToSentItems := false
	requestBody.SetSaveToSentItems(&saveToSentItems)

	// TODO: Correct
	return requestBody, nil
}

func toGraphRecipient(addresses ...*mail.Address) []graphmodels.Recipientable {
	var result = make([]graphmodels.Recipientable, len(addresses))

	for _, r := range addresses {
		graphRecipient := graphmodels.NewRecipient()
		graphAddress := graphmodels.NewEmailAddress()
		graphAddress.SetAddress(&r.Address)
		graphAddress.SetName(&r.Name)
		graphRecipient.SetEmailAddress(graphAddress)

		result = append(result, graphRecipient)
	}

	return result
}
