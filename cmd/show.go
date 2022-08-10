package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/avtools"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		//avtools.NewFFmpegCmd(input).Options(&flags).ShowMeta()
		f := avtools.NewMedia(input)
		f.AddFormat(flags.MetaFile)
		f.AddFormat(flags.CueFile)
		fmt.Printf("%+V\n", f.ListFormats())
		fmt.Printf("%+V\n", f.GetFormat("ini"))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.PersistentFlags().StringVarP(&flags.MetaFile, "metaFile", "m", "", "update ffmetadata")
	showCmd.PersistentFlags().StringVarP(&flags.CueFile, "cue", "c", "", "update chapter metadata")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
