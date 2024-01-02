package server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/smtp"
	"strings"
)

var _ = Describe("Server", func() {
	Context("sending mail", func() {
		It("should send mail", func() {
			auth := smtp.PlainAuth("", "user", "pass", "localhost")
			err := smtp.SendMail("localhost:5252", auth, "sender@localhost", []string{"recipient@localhost"}, createEmailMessage("hei", "this is a msg"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not send mail", func() {
			auth := smtp.PlainAuth("", "wronguser", "pass", "localhost")
			err := smtp.SendMail("localhost:5252", auth, "sender@localhost", []string{"recipient@localhost"}, createEmailMessage("hei", "this is a msg"))
			Expect(err).To(HaveOccurred())
		})
	})
})

func createEmailMessage(subject, body string) []byte {
	headers := []string{
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding: 7bit",
	}

	headersString := strings.Join(headers, "\r\n")
	return []byte(headersString + "\r\n\r\n" + body)
}
