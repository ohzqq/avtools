package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/meta"
	"github.com/spf13/cobra"
)

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "A brief description of your command",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		media := avtools.NewMedia(input)
		//fmt.Printf("media %+V\n", media.FFmeta.Chapters[0])
		//fmt.Println(c.String())
		//meta := meta.FFProbe(input)
		meta := meta.LoadIni(input)
		media.SetMeta(meta)
		fmt.Printf("meta %+V\n", media)
	},
}

func init() {
	rootCmd.AddCommand(probeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// probeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// probeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
