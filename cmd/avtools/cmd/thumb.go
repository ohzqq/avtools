package cmd

import (
	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var thumb media.Command

// thumbCmd represents the thumb command
var thumbCmd = &cobra.Command{
	Use:   "thumb",
	Short: "create a thumbnail from a video",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		tn := thumb.Thumbnail(input, outName)
		tn.Run()
	},
}

func init() {
	rootCmd.AddCommand(thumbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// thumbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// thumbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
