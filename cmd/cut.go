package cmd

import (
	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "cut files",
	Long:  `This can cut a file based either on provided timestamps or using a chapter number.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flag.Args.Input = args[0]
		u := tool.Cut()
		u.SetFlags(flag)
		c := u.Parse()
		c.RunBatch()
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)
	cutCmd.Flags().StringVarP(&flag.Args.Start, "start", "s", "", "start of clip")
	cutCmd.Flags().StringVarP(&flag.Args.End, "end", "e", "", "end of clip")
	cutCmd.Flags().IntVarP(&flag.Args.ChapNo, "num", "n", 0, "chapter number")
	cutCmd.MarkFlagsMutuallyExclusive("start", "num")
	cutCmd.MarkFlagsMutuallyExclusive("end", "num")
}
