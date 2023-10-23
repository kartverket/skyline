package cmd

import (
	"github.com/kartverket/skyline/pkg/config"
	"github.com/kartverket/skyline/pkg/server"
	"github.com/kartverket/skyline/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SMTP server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.NewServer(constructConfig())
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

	serveCmd.Flags().IntP("port", "p", 5252, "The SMTP port to bind to")
	serveCmd.Flags().IntP("metrics-port", "m", 5353, "The port to serve metrics at")
	serveCmd.Flags().String("hostname", h, "The SMTP hostname")

	serveCmd.Flags().String("sender-type", "msGraph", "Which underlying sending mechanism to use. Permitted values are 'msGraph' and 'dummy'")
	serveCmd.Flags().String("ms-tenant-id", "", "MS Graph API Tenant ID")
	serveCmd.Flags().String("ms-client-id", "", "MS Graph API Client ID")
	serveCmd.Flags().String("ms-client-secret", "", "MS Graph API Client Secret")
	serveCmd.Flags().String("ms-sender-id", "", "MS Graph API Object Id for Azure AD email sender")

	serveCmd.Flags().Bool("auth-enabled", true, "Whether basic auth is enabled")
	serveCmd.Flags().String("auth-username", "username", "Basic Auth username")
	serveCmd.Flags().String("auth-password", "password", "Basic Auth password")

	if err := viper.BindPFlags(serveCmd.Flags()); err != nil {
		slog.Error("could not bind cmd flags to viper config", "error", err)
		os.Exit(1)
	}
	_ = viper.BindPFlag("ms-graph-config.tenant-id", serveCmd.Flags().Lookup("ms-tenant-id"))
	_ = viper.BindPFlag("ms-graph-config.client-id", serveCmd.Flags().Lookup("ms-client-id"))
	_ = viper.BindPFlag("ms-graph-config.client-secret", serveCmd.Flags().Lookup("ms-client-secret"))
	_ = viper.BindPFlag("ms-graph-config.sender-user-id", serveCmd.Flags().Lookup("ms-sender-id"))

	_ = viper.BindPFlag("basic-auth-config.enabled", serveCmd.Flags().Lookup("auth-enabled"))
	_ = viper.BindPFlag("basic-auth-config.username", serveCmd.Flags().Lookup("auth-username"))
	_ = viper.BindPFlag("basic-auth-config.password", serveCmd.Flags().Lookup("auth-password"))
}

func constructConfig() *config.SkylineConfig {
	var cfg config.SkylineConfig
	_ = viper.Unmarshal(&cfg)

	if !cfg.SenderType.IsValid() {
		slog.Error("unknown sender type", "type", cfg.SenderType)
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
