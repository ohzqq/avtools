package cmd

import (
	"fmt"

	"github.com/ohzqq/avtools/tool"
	"github.com/ohzqq/avtools/tool/ffmpeg"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		tool.NewFFmpegCmd(input).Options(flags).ShowMeta()
		c := ffmpeg.New()
		c.SetInput(input)
		//c.AddCover(flags.CoverFile)
		//c.AddMeta(flags.MetaFile)
		eq := ffmpeg.NewFilter("eq")
		eq.Set("brightness", "1.0")
		eq.Set("saturation", "1.0")
		fps := ffmpeg.NewFilter("fps")
		fps.Set("60")
		c.AppendVideoFilter(eq)
		c.AppendVideoFilter(fps)
		c.AppendAudioParam("id3v2_version", "3")
		//f := avtools.NewMedia(input)
		//f.AddFormat(flags.MetaFile)
		//f.AddFormat(flags.CueFile)
		fmt.Printf("%+V\n", c.String())
		//fmt.Printf("%+V\n", f.GetFormat("ini"))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.PersistentFlags().StringVarP(&flags.MetaFile, "ffmeta", "f", "", "update ffmetadata")
	showCmd.PersistentFlags().StringVarP(&flags.CoverFile, "album-art", "a", "", "update album art")
	showCmd.PersistentFlags().StringVarP(&flags.CueFile, "cuesheet", "c", "", "update chapter metadata")
	showCmd.PersistentFlags().BoolVarP(&flags.MetaSwitch, "meta", "m", false, "print ffmeta")
	showCmd.PersistentFlags().BoolVarP(&flags.CueSwitch, "cue", "C", false, "print cue sheet")
	showCmd.PersistentFlags().BoolVarP(&flags.JsonSwitch, "json", "j", false, "print json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
