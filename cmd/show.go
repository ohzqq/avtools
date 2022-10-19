package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/media"
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
		flag.Args.Input = input

		medias(input)
	},
}

func medias(input string) {
	//c := media.NewMedia(input)
	c := media.NewCmd().Input(input).SetFlags(flag)
	c.Media.SetFFmeta(flag.Args.Meta)
	c.Media.SetCue(flag.Args.Cue)
	ff := c.NewFFmpegCmd()
	eq := ffmpeg.NewFilter("eq")
	eq.Set("brightness", "1.0")
	fps := ffmpeg.NewFilter("fps")
	fps.Set("60")
	ff.SetVideoCodec("libx264").
		AppendVideoFilter(eq).
		AppendVideoFilter(fps)

	ffcmd, err := ff.Build()
	if err != nil {
		log.Fatal(err)
	}
	c.AddCmd(ffcmd)
	fmt.Printf("%+V\n", c.Batch[0].String())
}

func ffdump(input string) {
	data := tool.EmbeddedMeta(input)
	ff := ffmeta.LoadJson(data)
	//ffmeta := ffmeta.LoadJson(data)
	ch := chap.NewChapters().FromCue(flag.Args.Cue)
	ff.SetChapters(ch)
	ff.Save()

	//for _, c := range u.Cmd.Batch {
	//c.Run()
	fmt.Printf("%+V\n", string(ff.Dump()))
	fmt.Printf("dur %+V\n", ff.Duration().String())
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
