package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/cue"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
)

var probe media.Command

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "show media info",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		//testIni(input)
		testCue(input)
		//metad := meta.FFProbe(input)
		//for _, ch := range metad.ChapterEntry {
		//fmt.Printf("probe meta %+V\n", ch.Start)
		//d := meta.ParseStamp(ch.Start, ch.Base)
		//f, err := dur.Parse(ch.Start)
		//if err != nil {
		//log.Fatal(err)
		//}
		//fmt.Printf("probe meta %+V\n", f.HHMMSS())
		//}
		//m := media.New(input)
		//for _, ch := range m.Chapters() {
		//  fmt.Printf("probe meta %+V\n", ch.Start.String())
		//}
		//fmt.Printf("meta %+V\n", m.Input.Abs)
		//fmt.Printf("meta %+V\n", len(m.Chapters()))
		//fmt.Printf("tags %+V\n", m.Chapters()[len(m.Chapters())-1].Tags)
		//meta := meta.LoadIni(input)
	},
}

func testCue(in string) {
	ff := cue.Load(in)
	fmt.Printf("%+V\n", ff)
	for _, ch := range ff.Chapters() {
		fmt.Printf("%s\n", ch.Start())
		fmt.Printf("%s\n", ch.End().Milliseconds())
	}
}

func testIni(in string) {
	ff := ffmeta.Load(in)
	fmt.Printf("%+V\n", ff)
	for _, ch := range ff.Chapters() {
		fmt.Printf("%s\n", ch.Start())
		fmt.Printf("%s\n", ch.End().Milliseconds())
	}
}

func init() {
	rootCmd.AddCommand(probeCmd)

	probeCmd.PersistentFlags().StringVarP(&probe.Flags.File.Meta, "meta", "m", "", "extract ffmeta")
	probeCmd.PersistentFlags().StringVarP(&probe.Flags.File.Cue, "cue", "c", "", "extract cue sheet")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// probeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// probeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
