package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var remove media.Command

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove elements from media",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		mCmd := remove.Remove(input)
		err := mCmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.PersistentFlags().BoolVarP(&remove.Bool.Meta, "meta", "m", false, "rm embedded meta")
	rmCmd.PersistentFlags().BoolVarP(&remove.Bool.Chapters, "chapters", "c", false, "rm chapters")
	rmCmd.PersistentFlags().BoolVarP(&remove.Bool.Cover, "album art", "a", false, "rm album art")
}
