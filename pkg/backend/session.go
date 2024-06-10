package backend

import (
	"context"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/kartverket/skyline/pkg/config"
	"github.com/kartverket/skyline/pkg/email"
	"github.com/kartverket/skyline/pkg/sender"
	"io"
	"log/slog"
)

// A Session is returned after EHLO.
type Session struct {
	auth      bool
	ctx       context.Context
	log       *slog.Logger
	sender    sender.Sender
	basicAuth *config.BasicAuthConfig
}

func (s *Session) AuthMechanisms() []string {
	if s.basicAuth.Enabled {
		return []string{sasl.Plain}
	}

	return []string{sasl.Anonymous}
}

func (s *Session) Auth(_ string) (sasl.Server, error) {
	if s.basicAuth.Enabled {
		return sasl.NewPlainServer(func(identity, username, password string) error {
			if len(identity) > 0 {
				s.log.Warn("user supplied identity but we don't support them", "identity", identity)
			}

			if !(username == s.basicAuth.Username && password == s.basicAuth.Password) {
				authenticationFailures.Inc()
				return fmt.Errorf("the password or the supplied username is incorrect")
				s.log.Warn("authentication failed", "username", username)
			}

			authenticationSuccesses.Inc()
			s.auth = true
			return nil
		}), nil
	}

	// No auth
	return sasl.NewAnonymousServer(func(trace string) error {
		s.log.Info("anonymous authentication", "trace", trace)
		return nil
	}), nil
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

	err = (s.sender).Send(s.ctx, msg)
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
	s.log.Debug("reset invoked")
}

func (s *Session) Logout() error {
	s.log.Debug("logout invoked")
	return nil
}
