package backend

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/kartverket/skyline/pkg/config"
	"github.com/kartverket/skyline/pkg/sender"
	"log/slog"
	"os"
)

// The Backend implements SMTP server methods.
// TODO: Configurable credentials
// TODO: Pluggable providers
type Backend struct {
	sender    *sender.Sender
	basicAuth *config.BasicAuthConfig
}

func NewBackend(cfg *config.SkylineConfig) *Backend {
	return &Backend{
		sender:    createSender(cfg),
		basicAuth: cfg.BasicAuthConfig,
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
	if !b.basicAuth.Enabled {
		slog.Warn("basic auth disabled, but validation called anyway")
		return true
	}

	return username == b.basicAuth.Username && password == b.basicAuth.Password
}

func createSender(cfg *config.SkylineConfig) *sender.Sender {
	var configuredSender sender.Sender

	switch cfg.SenderType {
	case config.MsGraph:
		s, err := sender.NewOffice365Sender(
			cfg.MsGraphConfig.TenantId,
			cfg.MsGraphConfig.ClientId,
			cfg.MsGraphConfig.ClientSecret,
			cfg.MsGraphConfig.SenderUserId,
		)
		if err != nil {
			slog.Error("could not construct sender", "error", err)
			os.Exit(1)
		}
		configuredSender = s
	case config.Dummy:
		slog.Warn("not implemented yet, exiting cleanly")
		os.Exit(0)
	default:
		slog.Error("unknown sender type", "type", cfg.SenderType)
		os.Exit(1)
	}

	return &configuredSender
}
