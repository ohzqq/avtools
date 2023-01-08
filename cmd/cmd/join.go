package cmd

import (
	"log"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

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

		c := media.Join(ext, dir)
		c.Run()
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
