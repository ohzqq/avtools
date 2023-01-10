package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

type splitFlags struct {
	cue    string
	ffmeta string
}

var split media.Command

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split on chapter markers",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		cmds := split.Split(input)
		for _, c := range cmds {
			err := c.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
	splitCmd.PersistentFlags().StringVarP(&split.Flags.File.Cue, "cue", "c", "", "split by cue sheet")
	splitCmd.PersistentFlags().StringVarP(&split.Flags.File.Meta, "meta", "m", "", "split by ffmetadata")
	splitCmd.MarkFlagsMutuallyExclusive("cue", "meta")
}
