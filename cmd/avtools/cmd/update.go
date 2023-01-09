package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/ff"
	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var update fmtStringFlags

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update metadata or cover art",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		m := media.New(input).Probe()
		in := ff.New()
		in.In(m.Input.Abs)
		if update.Meta != "" {
			in.Input.FFMeta(update.Meta)
		}
		in.Output.Set("c", "copy")
		//out := ffmpeg.Input("ffmeta.ini").Output("jlk", ffmpeg.KwArgs{"map_metadata": "1"})
		fmt.Printf("meta %+V\n", m.Input.Abs)
		fmt.Printf("args %+V\n", in.Compile().Args)
	},
}

func MapMetadata(file string, idx ...string) []string {
	label := "1"
	if len(idx) > 0 {
		label = idx[0]
	}
	input := []string{"-i", file, "-map_metadata", label}
	return input
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringVarP(&update.Meta, "meta", "m", "", "extract ffmeta")
	updateCmd.PersistentFlags().StringVarP(&update.Cue, "cue", "c", "", "extract cue sheet")
}
