package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "cut files",
	Long:  `This can cut a file based either on provided timestamps or using a chapter number.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(flags).Cut(flags.Start, flags.End, flags.ChapNo)
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)
	cutCmd.Flags().StringVarP(&flags.Start, "start", "s", "", "start of clip")
	cutCmd.Flags().StringVarP(&flags.End, "end", "e", "", "end of clip")
	cutCmd.Flags().IntVarP(&flags.ChapNo, "num", "n", 0, "chapter number")
	cutCmd.MarkFlagsMutuallyExclusive("start", "num")
	cutCmd.MarkFlagsMutuallyExclusive("end", "num")
}
