package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var join media.Command

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "join media files",
	Run: func(cmd *cobra.Command, args []string) {
		dir := "."
		var ext string
		switch len(args) {
		case 2:
			dir = args[1]
			fallthrough
		case 1:
			ext = args[0]
		default:
			log.Fatalf("wrong number of args")
		}

		ff, formats := media.Join(ext, dir)
		ff.Run()
		for format, c := range formats {
			if format == "ini" && join.Flags.Bool.Meta {
				c.Run()
			}
			if format == "cue" && join.Flags.Bool.Cue {
				c.Run()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
	joinCmd.PersistentFlags().BoolVarP(&join.Flags.Bool.Meta, "meta", "m", false, "extract ffmeta")
	joinCmd.PersistentFlags().BoolVarP(&join.Flags.Bool.Cue, "cue", "c", false, "extract cue sheet")
}
