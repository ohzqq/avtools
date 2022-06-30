package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the cover or metadata of a file",
	Long:  `Adds album art or updates the metadata (using ffmpeg's metadata format). Album art of aac requires Atomic Parsley`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(&flags).Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
