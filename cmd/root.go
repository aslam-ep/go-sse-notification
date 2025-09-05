package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "sse-notifications",
	Short: "SSE Notification Serivce",
	Long:  "A server-sent events notification system with Redis Pub/Sub and per-user streams.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Load .env if exists (optional)
	_ = godotenv.Load()
	viper.AutomaticEnv()
}
