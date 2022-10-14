package cmd

import (
	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract chapters, album art, or metadata",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tool.NewFFmpegCmd(args[0]).Options(flags).Extract()
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().BoolVarP(&flags.MetaSwitch, "meta", "m", false, "save ffmetadata to disk")
	extractCmd.PersistentFlags().BoolVarP(&flags.CueSwitch, "cue", "c", false, "save cue sheet to disk")
	extractCmd.PersistentFlags().BoolVarP(&flags.CoverSwitch, "album-art", "a", false, "save album art to disk")
}
