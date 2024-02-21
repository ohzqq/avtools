package cmd

import (
	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split a/v files",
	Long:  `split files by embedded chapters markers, an ffmpeg metadata file, or a cue sheet`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flag.Args.Input = args[0]
		u := tool.Split()
		u.ParseFlags(flag)
		c := u.Parse()
		c.RunBatch()
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
	splitCmd.PersistentFlags().StringVarP(&flag.Args.Cue, "cue", "c", "", "split by cue sheet")
	splitCmd.PersistentFlags().StringVarP(&flag.Args.Meta, "ffmeta", "f", "", "split by ffmetadata")
}
