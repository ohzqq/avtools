package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/avtools/av"
	"github.com/ohzqq/avtools/cue"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/media"
	"github.com/ohzqq/avtools/probe"
	"github.com/spf13/cobra"
)

var prob media.Command

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "show media info",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		//testIni(input)
		testMedia(input)
		//testCue(input)
		//testProbe(input)
		//metad := meta.FFProbe(input)
		//for _, ch := range metad.ChapterEntry {
		//fmt.Printf("probe meta %+V\n", ch.Start)
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

func testMedia(in string) {
	media := av.New(in)
	fmt.Printf("%+V\n", media)
	media.Cue("tmp/curse.cue")
	fmt.Printf("%+V\n", media)
	media.Probe()
	fmt.Printf("%+V\n", media)
}

func testProbe(in string) {
	ff, err := probe.Load(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+V\n", ff)
	//for _, ch := range ff.Chapters() {
	for _, ch := range ff.Chapters() {
		//fmt.Printf("%s\n", ch.StartTime)
		fmt.Printf("%s\n", ch.Start())
		fmt.Printf("%s\n", ch.End())
		//fmt.Printf("%s\n", ch.End().Milliseconds())
	}
}

func testCue(in string) {
	ff, err := cue.Load(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+V\n", ff)
	for _, ch := range ff.Chapters() {
		fmt.Printf("%s\n", ch.Start())
		fmt.Printf("%s\n", ch.End().Milliseconds())
	}
}

func testIni(in string) {
	ff, err := ffmeta.Load(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+V\n", ff)
	for _, ch := range ff.Chapters() {
		fmt.Printf("%s\n", ch.Start())
		fmt.Printf("%s\n", ch.End().Milliseconds())
	}
}

func init() {
	rootCmd.AddCommand(probeCmd)

	probeCmd.PersistentFlags().StringVarP(&prob.Flags.File.Meta, "meta", "m", "", "extract ffmeta")
	probeCmd.PersistentFlags().StringVarP(&prob.Flags.File.Cue, "cue", "c", "", "extract cue sheet")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// probeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// probeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
