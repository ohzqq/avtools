/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "A brief description of your command",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(&flags).Split()
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// splitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// splitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
