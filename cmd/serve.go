package cmd

import (
	"context"
	"github.com/kartverket/skyline/pkg/server"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var (
	port        uint
	metricsPort uint
	hostname    string
	debug       bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SMTP server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		s := server.NewServer(ctx, port, metricsPort, hostname, debug)
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
}
