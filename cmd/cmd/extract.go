package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var extract fmtBoolFlags

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract metadata or cover art",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		m := media.New(input).Probe()
		var cmds []media.Cmd
		if extract.Cover {
			ff := m.ExtractCover()
			cmds = append(cmds, ff)
		}
		if extract.Cue {
			c := m.SaveMetaFmt("cue")
			cmds = append(cmds, c)
		}
		if extract.Meta {
			c := m.SaveMetaFmt("ffmeta")
			cmds = append(cmds, c)
		}

		for _, c := range cmds {
			err := c.Run()
			if err != nil {
				log.Fatal(err)
			}
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
