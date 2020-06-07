package cmd

import (
	"github.com/atthavit/myblog/cli-covid19/covid19"
	"github.com/spf13/cobra"
)

var selectedType covid19.Type

func init() {
	rootCmd.AddCommand(plotCmd)

	for _, t := range covid19.Types {
		t := t
		plotCmd.AddCommand(&cobra.Command{
			Use: string(t),
			Run: func(cmd *cobra.Command, args []string) {
				selectedType = t
			},
		})
	}
}

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot a 30-day graph",
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return covid19.Plot(selectedType)
	},
}
