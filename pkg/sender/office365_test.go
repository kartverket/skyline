package sender

import (
	. "github.com/kartverket/skyline/pkg/email"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/mnako/letters"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/mail"
	"time"
)

var _ = Describe("Office365", func() {
	recipient := []*mail.Address{{Address: "skip@kartverket.no", Name: "SKIP"}}
	ccrecipient := []*mail.Address{{Address: "cc@kartverket.no", Name: "cc"}}
	bccrecipient := []*mail.Address{{Address: "bcc@kartverket.no", Name: "bcc"}}
	subject := "this is a subject"
	txtMsg := "This is a test message"
	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	htmlMsg := "<html><body>This is a test message</body></html>"
	skylineEmail := SkylineEmail{
		Email: letters.Email{
			Headers: letters.Headers{
				ContentType: letters.ContentTypeHeader{
					ContentType: "text/plain",
				},
				To:      recipient,
				Cc:      ccrecipient,
				Bcc:     bccrecipient,
				Date:    date,
				Subject: subject,
			},
			Text: txtMsg,
		},
	}
	Context("when mapping", func() {
		It("should map basic email to graph email", func() {
			graphMail, err := mapToGraphMail(&skylineEmail)
			Expect(err).To(BeNil())
			Expect(graphMail.GetMessage().GetSubject()).To(Equal(&subject))
			Expect(graphMail.GetMessage().GetToRecipients()).To(Equal(toGraphRecipient(recipient...)))
			Expect(graphMail.GetMessage().GetCcRecipients()).To(Equal(toGraphRecipient(ccrecipient...)))
			Expect(graphMail.GetMessage().GetBccRecipients()).To(Equal(toGraphRecipient(bccrecipient...)))
			Expect(graphMail.GetMessage().GetSentDateTime()).To(Equal(&date))
			contentType := graphmodels.TEXT_BODYTYPE
			Expect(graphMail.GetMessage().GetBody().GetContentType()).To(Equal(&contentType))
			Expect(graphMail.GetMessage().GetBody().GetContent()).To(Equal(&txtMsg))
		})
		It("should map html email to graph email", func() {
			skylineEmail.Headers.ContentType.ContentType = "text/html"
			skylineEmail.HTML = htmlMsg
			graphMail, err := mapToGraphMail(&skylineEmail)
			Expect(err).To(BeNil())
			contentType := graphmodels.HTML_BODYTYPE
			Expect(graphMail.GetMessage().GetBody().GetContentType()).To(Equal(&contentType))
			Expect(graphMail.GetMessage().GetBody().GetContent()).To(Equal(&htmlMsg))
		})
		It("should error if not valid type", func() {
			skylineEmail.Headers.ContentType.ContentType = "err"
			_, err := mapToGraphMail(&skylineEmail)
			Expect(err).ToNot(BeNil())
		})
		It("should map MIME email to graph email, html type", func() {
			skylineEmail.Headers.ContentType.ContentType = "multipart/alternative"
			skylineEmail.HTML = htmlMsg
			graphMail, err := mapToGraphMail(&skylineEmail)
			Expect(err).To(BeNil())
			contentType := graphmodels.HTML_BODYTYPE
			Expect(graphMail.GetMessage().GetBody().GetContentType()).To(Equal(&contentType))
			Expect(graphMail.GetMessage().GetBody().GetContent()).To(Equal(&htmlMsg))
		})
	})
})
