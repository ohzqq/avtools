package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var probe fmtStringFlags

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "A brief description of your command",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		m := media.New(input).Probe()
		switch {
		case cmd.Flags().Changed("meta"):
			m.LoadMeta(probe.Meta)
			cue := m.DumpIni()
			println(string(cue))
		case cmd.Flags().Changed("cue"):
			m.LoadMeta(probe.Cue)
			//cue := meta.DumpCueSheet(m.Input.Abs, m.Media)
		}
		//cut.Chapter(3)
		//cut.Start("00:01.000").End("00:02.999")
		//c := cut.Compile()
		//c.Run()
		fmt.Printf("meta %+V\n", m.Chapters())
		//fmt.Printf("meta %+V\n", m.Chapters[0].Start.HHMMSS())
		//fmt.Printf("meta %+V\n", m.Chapters[10].Start.String())
		//fmt.Printf("meta %+V\n", m.Chapters[10].Start.HHMMSS())
		//fmt.Printf("meta %+V\n", m.Chapters[10].Start.MMSS())
		//cue := m.DumpFFMeta()
		//cue.Compile().Run()
		//fmt.Printf("cue %+V\n", string(cue))
		//m.SetMeta(meta)
		//meta := meta.FFProbe(input)
		//meta := meta.LoadIni(input)
	},
}

func init() {
	rootCmd.AddCommand(probeCmd)

	probeCmd.PersistentFlags().StringVarP(&probe.Meta, "meta", "m", "", "extract ffmeta")
	probeCmd.PersistentFlags().StringVarP(&probe.Cue, "cue", "c", "", "extract cue sheet")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// probeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// probeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}