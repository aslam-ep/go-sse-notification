package cmd

import (
	"fmt"
	"log"

	"github.com/aslam-ep/go-sse-notification/config"
	"github.com/aslam-ep/go-sse-notification/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the SSE API server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		srv := server.NewServer(cfg)

		fmt.Printf("Starting server on %s....\n", cfg.Server.Address)
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
