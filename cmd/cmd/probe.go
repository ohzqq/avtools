package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/media"
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
		//m := media.New(input).LoadCue(input)
		//m := media.New(input).LoadIni(input)
		m := media.New(input).Probe()
		fmt.Printf("meta %+V\n", m.Filename)
		c := m.ExtractCover()
		c.Compile().Run()
		//m.SetMeta(meta)
		//meta := meta.FFProbe(input)
		//meta := meta.LoadIni(input)
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
