package cmd

import (
	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract chapters, album art, or metadata",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//tool.NewFFmpegCmd(args[0]).Options(flags).Extract()
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().BoolVarP(&flag.Bool.Meta, "meta", "m", false, "print ffmeta")
	extractCmd.PersistentFlags().BoolVarP(&flag.Bool.Cue, "cue", "C", false, "print cue sheet")
	extractCmd.PersistentFlags().BoolVarP(&flag.Bool.Json, "json", "j", false, "print json")
}
