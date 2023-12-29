package server

import (
	"context"
	"fmt"
	"github.com/emersion/go-smtp"
	skybackend "github.com/kartverket/skyline/pkg/backend"
	"github.com/kartverket/skyline/pkg/config"
	"github.com/kartverket/skyline/pkg/email"
	skysender "github.com/kartverket/skyline/pkg/sender"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

type mockSender struct{}

func (m *mockSender) Send(ctx context.Context, email *email.SkylineEmail) error {
	return nil
}

var skyServer *SkylineServer

var _ = BeforeSuite(func() {
	cfg := config.SkylineConfig{
		Hostname:    "localhost",
		Port:        5252,
		MetricsPort: 5353,
		SenderType:  0,
		BasicAuthConfig: &config.BasicAuthConfig{
			Username: "user",
			Password: "pass",
			Enabled:  true,
		},
		MsGraphConfig: &config.MsGraphConfig{
			TenantId:     "tenant",
			ClientId:     "client",
			SenderUserId: "1",
			ClientSecret: "1",
		},
	}
	var sender skysender.Sender
	sender = &mockSender{}

	backend := &skybackend.Backend{
		Sender:    sender,
		BasicAuth: cfg.BasicAuthConfig,
	}
	smtpServer := smtp.NewServer(backend)

	smtpServer.Addr = fmt.Sprintf(":%d", cfg.Port)
	smtpServer.Domain = cfg.Hostname
	smtpServer.ReadTimeout = 10 * time.Second
	smtpServer.WriteTimeout = 10 * time.Second
	smtpServer.MaxMessageBytes = 1024 * 1024
	smtpServer.MaxRecipients = 50
	smtpServer.AllowInsecureAuth = true
	smtpServer.ErrorLog = logAdapter(ctx)

	skyServer = &SkylineServer{
		smtp:    smtpServer,
		ctx:     ctx,
		metrics: metricsServer(cfg.MetricsPort),
	}
	Expect(skyServer).ToNot(BeNil())
	go func() {
		skyServer.Serve()
	}()
})

var _ = AfterSuite(func() {
	skyServer.ctx.Done()
})
