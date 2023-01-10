package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var update media.Command

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update metadata or cover art",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		//m := media.Update(input, update.Meta, update.Cue)
		m := update.Update(input)
		//switch {
		//case update.Meta != "":
		//  m.LoadMeta(update.Meta)
		//  //in.Input.FFMeta(update.Meta)
		//case update.Cue != "":
		//  m.LoadMeta(update.Cue)
		//}

		err := m.Run()
		if err != nil {
			log.Fatal(err)
		}
		//out := ffmpeg.Input("ffmeta.ini").Output("jlk", ffmpeg.KwArgs{"map_metadata": "1"})
		//fmt.Printf("args %+V\n", in.Compile().Args)
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
	updateCmd.PersistentFlags().StringVarP(&update.Flags.File.Meta, "meta", "m", "", "extract ffmeta")
	updateCmd.PersistentFlags().StringVarP(&update.Flags.File.Cue, "cue", "c", "", "extract cue sheet")
}
