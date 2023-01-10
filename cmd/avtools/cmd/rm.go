package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var remove media.Command

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove elements from media",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rm called")
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.PersistentFlags().BoolVarP(&extract.Bool.Meta, "meta", "m", false, "extract ffmeta")
	rmCmd.PersistentFlags().BoolVarP(&extract.Bool.Cue, "cue", "c", false, "extract cue sheet")
	rmCmd.PersistentFlags().BoolVarP(&extract.Bool.Cover, "album art", "a", false, "extract album art")
}
