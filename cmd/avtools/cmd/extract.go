package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var extract media.Command

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract metadata or cover art",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		cmds := extract.Extract(input)

		for _, c := range cmds {
			err := c.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	// flags
	extractCmd.PersistentFlags().BoolVarP(&extract.Bool.Meta, "meta", "m", false, "extract ffmeta")
	extractCmd.PersistentFlags().BoolVarP(&extract.Bool.Cue, "cue", "c", false, "extract cue sheet")
	extractCmd.PersistentFlags().BoolVarP(&extract.Bool.Cover, "album art", "a", false, "extract album art")
}
