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
	rmCmd.PersistentFlags().BoolVarP(&flags.MetaSwitch, "meta", "m", false, "delete all embedded metadata")
	rmCmd.PersistentFlags().BoolVarP(&flags.CoverSwitch, "album-art", "a", false, "remove embedded album art")
	rmCmd.PersistentFlags().BoolVarP(&flags.ChapSwitch, "chapters", "c", false, "remove embedded album art")
}
