package cmd

import (
	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove/delete album art, metadata, or chapters",
	Long:  `Use this to remove metadata, art, or chapters from a/v files`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flag.Args.Input = args[0]
		e := tool.Rm()
		e.SetFlags(flag)
		c := e.Parse()
		c.RunBatch()
		//fmt.Printf("%+V\n", c)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().BoolVarP(&flag.Bool.Meta, "meta", "m", false, "print ffmeta")
	rmCmd.PersistentFlags().BoolVarP(&flag.Bool.Cue, "cue", "c", false, "print cue sheet")
	rmCmd.PersistentFlags().BoolVarP(&flag.Bool.Cover, "album art", "a", false, "print cue sheet")
}
