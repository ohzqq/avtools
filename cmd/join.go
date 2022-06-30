package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "concantenate a/v files",
	Long: `This uses ffmpeg's concat demuxer, which requires all files to have the same streams. This particular script requires that they all have the same container/extension. Best used on files that were split from a single source.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd("").Options(&flags).Join(args[0])
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
