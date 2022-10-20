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
		flag.Args.Input = args[0]
		e := tool.Extract()
		e.SetFlags(flag)
		c := e.Parse()
		c.RunBatch()
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().BoolVarP(&flag.Bool.Meta, "meta", "m", false, "print ffmeta")
	extractCmd.PersistentFlags().BoolVarP(&flag.Bool.Cue, "cue", "c", false, "print cue sheet")
	extractCmd.PersistentFlags().BoolVarP(&flag.Bool.Cover, "album art", "a", false, "print cue sheet")
}
