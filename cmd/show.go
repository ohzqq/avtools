/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/ohzqq/avtools/avtools"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(&flags).ShowMeta()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.PersistentFlags().StringVarP(&flags.MetaFile, "metaFile", "m", "", "update ffmetadata")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
