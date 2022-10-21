package cmd

import (
	"github.com/ohzqq/avtools/tool"
	"github.com/spf13/cobra"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "concantenate a/v files",
	Long:  `This uses ffmpeg's concat demuxer, which requires all files to have the same streams. This particular script requires that they all have the same container/extension. Best used on files that were split from a single source.`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		ext := args[0]
		dir := "."
		if len(args) == 2 {
			dir = args[1]
		}
		jcmd := tool.Join(ext, dir)
		jcmd.ParseFlags(flag)
		c := jcmd.Parse()
		c.RunBatch()
		//fmt.Printf("%+V\n", c.String())
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
