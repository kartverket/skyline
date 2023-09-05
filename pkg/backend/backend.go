package backend

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/kartverket/skyline/pkg/sender"
	"log/slog"
	"os"
)

// The Backend implements SMTP server methods.
// TODO: Configurable credentials
// TODO: Pluggable providers
type Backend struct {
	sender   *sender.Sender
	username string
	password string
}

func NewBackend(username string, password string) *Backend {
	// TODO: Add config and dummy sender
	s, err := sender.NewOffice365Sender("", "", "", "")
	if err != nil {
		slog.Error("could not construct sender", "error", err)
		os.Exit(1)
	}

	return &Backend{
		sender:   &s,
		username: username,
		password: password,
	}
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	// TODO: Plug in OpenTelemetry
	u, _ := uuid.NewUUID()
	ctx := context.WithValue(context.Background(), "trace-id", u.String())
	logger := slog.Default().With("trace-id", u.String())

	return &Session{
		ctx:                 ctx,
		log:                 logger,
		sender:              b.sender,
		validateCredentials: b.checkCredentials,
	}, nil
}

func (b *Backend) checkCredentials(username string, password string) bool {
	return username == b.username && password == b.password
}
