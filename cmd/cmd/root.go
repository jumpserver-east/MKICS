package cmd

import (
	"MKICS/backend/server"

	"github.com/spf13/cobra"
)

var configFile string

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/opt/mkics/conf/config.yaml", "Config file path")
}

var RootCmd = &cobra.Command{
	Use:   "MKICS",
	Short: "MKICS",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start(configFile)
		return nil
	},
}
