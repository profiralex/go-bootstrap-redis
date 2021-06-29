package cmd

import (
	"github.com/profiralex/go-bootstrap-redis/pkg/config"
	"github.com/profiralex/go-bootstrap-redis/pkg/db"
	"github.com/profiralex/go-bootstrap-redis/pkg/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "command to serve the app over http",
	Long:  "command to serve the app over http",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		db.Init(cfg)
		server.Serve(cfg)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
