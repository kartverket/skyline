package cmd

import (
	"context"
	"github.com/kartverket/skyline/pkg/config"
	"github.com/kartverket/skyline/pkg/server"
	"github.com/kartverket/skyline/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var (
	port        uint
	metricsPort uint
	hostname    string
	debug       bool
	senderType  string

	msGraphTenantId     string
	msGraphClientId     string
	msGraphClientSecret string
	msGraphSenderUserId string

	baEnabled  bool
	baUsername string
	baPassword string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SMTP server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		s := server.NewServer(ctx, constructConfig())
		s.Serve()
	},
}

func init() {
	h, err := os.Hostname()
	if err != nil {
		slog.Error("could not get hostname", "error", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().UintVarP(&port, "port", "p", 5252, "The SMTP port to bind to")
	serveCmd.Flags().UintVarP(&metricsPort, "metrics-port", "m", 5353, "The port to serve metrics at")
	serveCmd.Flags().StringVar(&hostname, "hostname", h, "The SMTP hostname")
	serveCmd.Flags().BoolVar(&debug, "debug", false, "Whether to enable debug logging")

	serveCmd.Flags().StringVar(&senderType, "sender-type", "msGraph", "Which underlying sending mechanism to use. Permitted values are 'msGraph' and 'dummy'")
	serveCmd.Flags().StringVar(&msGraphTenantId, "ms-tenant-id", "", "MS Graph API Tenant ID")
	serveCmd.Flags().StringVar(&msGraphClientId, "ms-client-id", "", "MS Graph API Client ID")
	serveCmd.Flags().StringVar(&msGraphClientSecret, "ms-client-secret", "", "MS Graph API Client Secret")
	serveCmd.Flags().StringVar(&msGraphSenderUserId, "ms-sender-id", "", "MS Graph API Object Id for Azure AD email sender")

	serveCmd.Flags().BoolVar(&baEnabled, "auth-enabled", true, "Whether basic auth is enabled")
	serveCmd.Flags().StringVar(&baUsername, "auth-username", "username", "Basic Auth username")
	serveCmd.Flags().StringVar(&baPassword, "auth-password", "password", "Basic Auth password")

	if err := viper.BindPFlags(serveCmd.Flags()); err != nil {
		slog.Error("could not bind cmd flags to viper config", "error", err)
		os.Exit(1)
	}
}

func constructConfig() *config.SkylineConfig {
	var cfg config.SkylineConfig
	_ = viper.Unmarshal(&cfg)

	switch t := cfg.SenderType; t {
	case config.MsGraph:
	case config.Dummy:
	default:
		slog.Error("unknown sender type", "type", t)
		os.Exit(1)
	}

	if cfg.SenderType == config.MsGraph {
		if cfg.MsGraphConfig == nil || util.AnyEmpty(cfg.MsGraphConfig.ClientId, cfg.MsGraphConfig.TenantId, cfg.MsGraphConfig.SenderUserId, cfg.MsGraphConfig.ClientSecret) {
			slog.Error("sender is configured as MsGraph but some of the required configuration properties is empty")
			os.Exit(1)
		}
	}

	return &cfg
}
