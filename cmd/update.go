package cmd

import (
	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the cover or metadata of a file",
	Long:  `Adds album art or updates the metadata (using ffmpeg's metadata format). Album art of aac requires Atomic Parsley`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flag.Args.Input = args[0]
		u := tool.NewUpdateCmd()
		u.SetFlags(flag)
		c := u.Parse()
		c.RunBatch()
		//for _, c := range u.Cmd.Batch {
		//c.Run()
		//}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringVarP(&flag.Args.Meta, "ffmeta", "m", "", "update ffmetadata")
	updateCmd.PersistentFlags().StringVarP(&flag.Args.Cover, "album-art", "a", "", "update album art")
	updateCmd.PersistentFlags().StringVarP(&flag.Args.Cue, "cuesheet", "c", "", "update chapter metadata")
}
