package cmd

import (
	"github.com/atthavit/myblog/cli-covid19/covid19"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(todayCmd)
}

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Print today's stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		return covid19.PrintToday()
	},
}
