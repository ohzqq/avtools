package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

type cutFlags struct {
	start string
	end   string
	chap  int
}

var cut cutFlags

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "A brief description of your command",
	Long:  `losslessly cut a clip`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		c := media.Cut(input)
		if cmd.Flags().Changed("ss") {
			c.Start(cut.start)
		}
		if cmd.Flags().Changed("to") {
			c.End(cut.end)
		}
		if cmd.Flags().Changed("num") {
			c.Chapter(cut.chap)
		}
		err := c.Compile().Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)
	cutCmd.Flags().StringVarP(&cut.start, "ss", "s", "", "start of clip")
	cutCmd.Flags().StringVarP(&cut.end, "to", "e", "", "end of clip")
	cutCmd.Flags().IntVarP(&cut.chap, "num", "n", 0, "chapter number")
	cutCmd.MarkFlagsMutuallyExclusive("ss", "num")
	cutCmd.MarkFlagsMutuallyExclusive("to", "num")
}
