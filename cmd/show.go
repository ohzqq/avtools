package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		//tool.NewFFmpegCmd(input).Options(flags).ShowMeta()
		//c := ffmpeg.New()
		flag.Args.Input = input
		u := tool.NewUpdateCmd() //.SetFlags(flag)
		u.Cmd = tool.NewerCmd().SetFlags(flag)
		u.ParseArgs()
		//c.Input(input).CA("libfdk_aac").AppendPreInput("y")
		//eq := ffmpeg.NewFilter("eq")
		//eq.Set("brightness", "1.0")
		//eq.Set("saturation", "1.0")
		//fps := ffmpeg.NewFilter("fps")
		//fps.Set("60")
		//c.Filter(eq).VF(fps).AppendAudioParam("id3v2_version", "3").Output("tep.m4b")
		//f := avtools.NewMedia(input)
		//f.AddFormat(flags.MetaFile)
		//f.AddFormat(flags.CueFile)
		//ff, err := c.Build()
		//if err != nil {
		//log.Fatal(err)
		//}
		for _, c := range u.Cmd.Batch {
			c.Run()
			fmt.Printf("%+V\n", c.String())
		}
		//fmt.Printf("%+V\n", f.GetFormat("ini"))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.PersistentFlags().StringVarP(&flag.Args.Output, "output", "o", "", "update ffmetadata")
	showCmd.PersistentFlags().StringVarP(&flag.Args.Meta, "ffmeta", "f", "", "update ffmetadata")
	showCmd.PersistentFlags().StringVarP(&flag.Args.Cover, "album-art", "a", "", "update album art")
	showCmd.PersistentFlags().StringVarP(&flag.Args.Cue, "cuesheet", "c", "", "update chapter metadata")
	showCmd.PersistentFlags().BoolVarP(&flag.Bool.Meta, "meta", "m", false, "print ffmeta")
	showCmd.PersistentFlags().BoolVarP(&flag.Bool.Cue, "cue", "C", false, "print cue sheet")
	showCmd.PersistentFlags().BoolVarP(&flag.Bool.Json, "json", "j", false, "print json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
