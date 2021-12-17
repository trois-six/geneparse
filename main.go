package main

import (
	"os"

	"github.com/Trois-Six/geneparse/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "geneparse",
		Short: "A tool to download, extract and parse Geneanet bases.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand((&cmd.DownloadAndExtractCmd{}).Command())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
