package backend

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/kartverket/skyline/pkg/config"
	skysender "github.com/kartverket/skyline/pkg/sender"
	"log/slog"
	"os"
)

// The Backend implements SMTP server methods.
type Backend struct {
	Sender    skysender.Sender
	BasicAuth *config.BasicAuthConfig
}

func NewBackend(cfg *config.SkylineConfig) *Backend {
	return &Backend{
		Sender:    createSender(cfg),
		BasicAuth: cfg.BasicAuthConfig,
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
		sender:              b.Sender,
		validateCredentials: b.checkCredentials,
	}, nil
}

func (b *Backend) checkCredentials(username string, password string) bool {
	if !b.BasicAuth.Enabled {
		slog.Warn("basic auth disabled, but validation called anyway")
		return true
	}

	return username == b.BasicAuth.Username && password == b.BasicAuth.Password
}

func createSender(cfg *config.SkylineConfig) skysender.Sender {
	var configuredSender skysender.Sender

	switch cfg.SenderType {
	case config.MsGraph:
		sender, err := skysender.NewOffice365Sender(
			cfg.MsGraphConfig.TenantId,
			cfg.MsGraphConfig.ClientId,
			cfg.MsGraphConfig.ClientSecret,
			cfg.MsGraphConfig.SenderUserId,
		)
		if err != nil {
			slog.Error("could not construct sender", "error", err)
			os.Exit(1)
		}
		configuredSender = sender
	case config.Dummy:
		slog.Warn("not implemented yet, exiting cleanly")
		os.Exit(0)
	default:
		slog.Error("unknown sender type", "type", cfg.SenderType)
		os.Exit(1)
	}

	return configuredSender
}
