package cmd

import (
	"fmt"
	"os"

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
		c := tool.NewCmd()
		c.ParseFlags(flag)
		fmt.Println(string(c.Media.Meta.Chapters.ToJson()))
		if cmd.Flags().Changed("cue") {
			c.Media.Meta.Chapters.ToCue().Write(os.Stdout)
		}
		if cmd.Flags().Changed("meta") {
			c.Media.Meta.Write(os.Stdout)
		}
		if cmd.Flags().Changed("json") {
			data := c.Media.Meta.DumpJson()
			fmt.Println(string(data))
		}
		//ff := media.Cfg().Profiles["gif"].FFmpegCmd()
		//ff.Input(input)

		//fmt.Printf("%+V\n", media.Cfg().ListProfiles())
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.PersistentFlags().BoolVarP(&flag.Bool.Meta, "meta", "m", false, "print ffmeta")
	showCmd.PersistentFlags().BoolVarP(&flag.Bool.Cue, "cue", "c", false, "print cue sheet")
	showCmd.PersistentFlags().BoolVarP(&flag.Bool.Json, "json", "j", false, "print json")
}
