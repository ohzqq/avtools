package cmd

import (
	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the cover or metadata of a file",
	Long:  `Adds album art or updates the metadata (using ffmpeg's metadata format). Album art of aac requires Atomic Parsley`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if flags.CueFile != "" {
			cue := avtools.NewFileFormat(flags.CueFile)
			cue.ConvertTo("cue").Print()
		}
		if flags.MetaFile != "" {
			cue := avtools.NewFileFormat(flags.MetaFile)
			cue.ConvertTo("json").Print()
		}
		avtools.NewFFmpegCmd(args[0]).Options(&flags).Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringVarP(&flags.CoverFile, "artFile", "a", "", "update album art")
	updateCmd.PersistentFlags().StringVarP(&flags.MetaFile, "metaFile", "m", "", "update ffmetadata")
	updateCmd.PersistentFlags().StringVarP(&flags.CueFile, "cue", "c", "", "update chapter metadata")
}
