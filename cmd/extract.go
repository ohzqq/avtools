package cmd

import (
	//"fmt"

	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract chapters, album art, or metadata",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(&flags).Extract()
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
