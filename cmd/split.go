package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split a/v files",
	Long:  `split files by embedded chapters markers, an ffmpeg metadata file, or a cue sheet`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		avtools.NewFFmpegCmd(args[0]).Options(&flags).Split()
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
	splitCmd.PersistentFlags().StringVarP(&flags.CueFile, "cuesheet", "c", "", "split file with cue sheet")
	splitCmd.PersistentFlags().StringVarP(&flags.MetaFile, "metaFile", "m", "", "split file with ffmetadata file")
}
