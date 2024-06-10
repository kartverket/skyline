package sender

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/kartverket/skyline/pkg/email"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/pkg/errors"
	"log/slog"
	"net/mail"
	"strings"
)

type office365sender struct {
	graphClient  *msgraphsdk.GraphServiceClient
	senderUserId string
}

func NewOffice365Sender(tenantId string, clientId string, clientSecret string, senderUserId string) (Sender, error) {
	slog.Info("Creating new client", "tenant-id", tenantId, "client-id", clientId)

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, err
	}

	c, err := msgraphsdk.NewGraphServiceClientWithCredentials(credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}

	return &office365sender{c, senderUserId}, nil
}

func (s *office365sender) Send(ctx context.Context, email *email.SkylineEmail) error {
	payload, err := mapToGraphMail(email)
	if err != nil {
		return errors.Wrap(err, "could not convert received message to fit within Microsoft Graph API models")
	}

	err = s.graphClient.Users().
		ByUserId(s.senderUserId).
		SendMail().Post(ctx, payload, nil)

	return unwrapODataError(err)
}

func mapToGraphMail(email *email.SkylineEmail) (*graphusers.ItemSendmailSendMailPostRequestBody, error) {
	// Outer payload
	requestBody := graphusers.NewItemSendmailSendMailPostRequestBody()
	message := graphmodels.NewMessage()

	// Subject
	message.SetSubject(&email.Headers.Subject)

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
	} else if email.IsHTML() || email.IsMultiPartAlternative() {
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

	return requestBody, nil
}

func toGraphRecipient(addresses ...*mail.Address) []graphmodels.Recipientable {
	var result = make([]graphmodels.Recipientable, 0)

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

// Microsoft seem unable to return proper errors, so we'll have to hack around it
func unwrapODataError(err error) error {
	if err == nil {
		return nil
	}

	var e *odataerrors.ODataError
	if errors.As(err, &e) {
		var sb strings.Builder
		sb.WriteString("[MS Graph API] ")
		sb.WriteString(e.Error())
		if details := e.GetErrorEscaped(); details != nil {
			sb.WriteString(fmt.Sprintf(" (code=%s, message='%s')", *details.GetCode(), *details.GetMessage()))
		}

		return errors.New(sb.String())
	}

	return err
}
