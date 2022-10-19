package cmd

import (
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove/delete album art, metadata, or chapters",
	Long:  `Use this to remove metadata, art, or chapters from a/v files`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//tool.NewFFmpegCmd(args[0]).Options(flags).Remove()
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().BoolVarP(&flag.Bool.Meta, "meta", "m", false, "print ffmeta")
	rmCmd.PersistentFlags().BoolVarP(&flag.Bool.Cue, "cue", "C", false, "print cue sheet")
	rmCmd.PersistentFlags().BoolVarP(&flag.Bool.Json, "json", "j", false, "print json")
}
