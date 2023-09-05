package backend

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/kartverket/skyline/pkg/email"
	"github.com/kartverket/skyline/pkg/sender"
	"io"
	"log/slog"
)

// TODO: Support disabling of username/password
// A Session is returned after EHLO.
type Session struct {
	auth                bool
	ctx                 context.Context
	log                 *slog.Logger
	sender              *sender.Sender
	validateCredentials func(string, string) bool
}

func (s *Session) AuthPlain(username, password string) error {
	if !s.validateCredentials(username, password) {
		authenticationFailures.Inc()
		s.log.Debug("Session authentication failed")
		return smtp.ErrAuthFailed
	}
	s.auth = true
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.auth {
		authenticationFailures.Inc()
		return smtp.ErrAuthRequired
	}
	s.log.Debug("Sender set", "address", from)
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.auth {
		authenticationFailures.Inc()
		return smtp.ErrAuthRequired
	}
	s.log.Debug("Recipient set", "address", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if !s.auth {
		authenticationFailures.Inc()
		return smtp.ErrAuthRequired
	}

	msg, err := email.Parse(r)
	if err != nil {
		emailParseFailures.Inc()
		return err
	}

	emailsProcessed.Inc()
	s.log.Debug("received and parsed email", "mail", msg)

	err = (*s.sender).Send(s.ctx, msg)
	if err != nil {
		s.log.Warn("could not send email", "error", err)
		emailsFailed.Inc()
		return err
	}

	s.log.Info("email sent OK")
	emailsSucceeded.Inc()
	return nil
}

func (s *Session) Reset() {
}

func (s *Session) Logout() error {
	return nil
}
