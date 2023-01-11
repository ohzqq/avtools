package cmd

import (
	"log"
	"os"

	"github.com/ohzqq/avtools/cmd/yygif/gif"
	"github.com/spf13/cobra"
)

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "A brief description of your command",
	//Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if MetaExists("metadata-default.yml") {
			gifMeta := gif.ReadMeta("metadata-default.yml")
			ini := gifMeta.DumpIni()
			file, err := os.Create("gif-meta.ini")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			_, err = file.Write(ini)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)
}
