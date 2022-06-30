/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "A brief description of your command",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(&flags).Cut(flags.Start, flags.End, flags.ChapNo)
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)
	cutCmd.Flags().StringVarP(&flags.Start, "start", "s", "", "start of clip")
	cutCmd.Flags().StringVarP(&flags.End, "end", "e", "", "end of clip")
	cutCmd.Flags().IntVarP(&flags.ChapNo, "num", "n", 0, "chapter number")
	cutCmd.MarkFlagsMutuallyExclusive("start", "num")
	cutCmd.MarkFlagsMutuallyExclusive("end", "num")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
