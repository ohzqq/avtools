package cmd

import (
	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

type splitFlags struct {
	cue    string
	ffmeta string
}

var split splitFlags

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split on chapter markers",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		cut := media.Cut(input)

		if cmd.Flags().Changed("cue") {
			cut.SetMeta(split.cue)
		}

		if cmd.Flags().Changed("meta") {
			cut.SetMeta(split.ffmeta)
		}

		cut.AllChapters()
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
	splitCmd.PersistentFlags().StringVarP(&split.cue, "cue", "c", "", "split by cue sheet")
	splitCmd.PersistentFlags().StringVarP(&split.ffmeta, "meta", "m", "", "split by ffmetadata")
	splitCmd.MarkFlagsMutuallyExclusive("cue", "meta")
}
