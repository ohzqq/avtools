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
		var gifMeta gif.Meta
		if !cmd.Flags().Changed("meta") {
			if MetaExists("metadata-default.yml") {
				gifMeta = gif.ReadMeta("metadata-default.yml")
				ini := gifMeta.DumpIni()
				println(string(ini))
			}
			if len(args) > 0 {
				arg := strings.Split(args[0], ",")
				if len(arg) != 2 {
					log.Fatalf("needs two args")
				}
				clip := gifMeta.GetClip(arg[0], arg[1])
				ff := ParseFlags(cmd, clip)
				ff.Compile()
				err := ff.Run()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				c := gifMeta.MkGifs()
				for _, clip := range c {
					clip.Compile()
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
