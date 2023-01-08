package cmd

import (
	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

type extractFlags struct {
	Meta  bool
	Cue   bool
	Cover bool
}

var extract extractFlags

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract metadata or cover art",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		m := media.New(input).Probe()
		if extract.Cover {
			m.ExtractCover()
		}
		if extract.Cue {
			m.SaveMetaFmt("cue")
		}
		if extract.Meta {
			m.SaveMetaFmt("ffmeta")
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	// flags
	extractCmd.PersistentFlags().BoolVarP(&extract.Meta, "meta", "m", false, "extract ffmeta")
	extractCmd.PersistentFlags().BoolVarP(&extract.Cue, "cue", "c", false, "extract cue sheet")
	extractCmd.PersistentFlags().BoolVarP(&extract.Cover, "album art", "a", false, "extract album art")
}
