package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var cut media.Command
var start string
var end string
var chap int

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "lossly cut clips",
	Long:  `losslessly cut a clip`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]

		var cutCmd media.Cmd

		if cmd.Flags().Changed("ss") || cmd.Flags().Changed("to") {
			cutCmd = cut.CutStamp(input, start, end)
		}

		if cmd.Flags().Changed("num") {
			cutCmd = cut.CutChapter(input, chap)
		}

		err := cutCmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)
	cutCmd.Flags().StringVarP(&start, "ss", "s", "", "start of clip")
	cutCmd.Flags().StringVarP(&end, "to", "e", "", "end of clip")
	cutCmd.Flags().IntVarP(&chap, "num", "n", 0, "chapter number")
	cutCmd.MarkFlagsMutuallyExclusive("ss", "num")
	cutCmd.MarkFlagsMutuallyExclusive("to", "num")
}
