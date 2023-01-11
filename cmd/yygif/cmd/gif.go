package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/ohzqq/avtools/cmd/yygif/gif"
	"github.com/spf13/cobra"
)

// gifCmd represents the gif command
var gifCmd = &cobra.Command{
	Use:   "gif",
	Short: "make gifs",
	Run: func(cmd *cobra.Command, args []string) {
		var meta gif.Meta
		if !cmd.Flags().Changed("meta") {
			if MetaExists(metadata) {
				meta = gif.ReadMeta(metadata)
			}
			if len(args) > 0 {
				arg := strings.Split(args[0], ",")
				if len(arg) != 2 {
					log.Fatalf("needs two args")
				}
				clip := meta.GetClip(arg[0], arg[1])
				ff := ParseFlags(cmd, clip)
				err := ff.Run()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				c := meta.MkGifs()
				for _, clip := range c {
					err := clip.Run()
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	},
}

func MetaExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	rootCmd.AddCommand(gifCmd)
}
