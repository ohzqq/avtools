/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"

	//"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: ``,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("%+v\n", avCmd.Show(args[0], args[1]))

	},
}

//func(cmd *Cmd) Show(action string) *Cmd {
//  media := NewMedia(cmd.input)
//  switch action {
//  case "json":
//    media.JsonMeta().Print()
//    //fmt.Printf("%+V\n", string(cmd.Media.GetJsonMeta()))
//  case "flags":
//    //fmt.Printf("%+v\n", cmd.Flags)
//  case "args":
//  case "chaps":
//    ch, err := cmd.getChapters()
//    if err != nil {
//      log.Fatal(err)
//    }
//    fmt.Printf("%+v\n", ch)
//  case "meta":
//    m := media.JsonMeta().Unmarshal()
//    fmt.Printf("%+V\n", m.Meta)
//  case "cmd":
//    //m := NewMedia(input).JsonMeta().Unmarshal()
//    //fmt.Printf("%+V\n", m.Meta)
//    //cmd.ffmpeg = true
//    //cmd.ffprobe = true
//    //fmt.Printf("%+v\n", Cfg().GetProfile(cmd.Flags.Profile))
//    //fmt.Printf("%+v\n", cmd.exec.String())
//  default:
//    fmt.Printf("%+v\n", cmd)
//  }
//  return cmd
//}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
