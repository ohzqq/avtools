package avtools

import (
	//"path/filepath"
	//"log"
	"fmt"
	//"os"
	//"strings"
	"strconv"
	"os/exec"

	//"github.com/alessio/shellescape"
)
var _ = fmt.Printf

type Args struct {
	Flags
	PreInput mapArgs
	PostInput mapArgs
	VideoCodec string
	VideoParams mapArgs
	VideoFilters stringArgs
	AudioCodec string
	AudioParams mapArgs
	AudioFilters stringArgs
	FilterComplex stringArgs
	MiscParams stringArgs
	LogLevel string
	Name string
	Padding string
	Ext string
	num int
	pretty bool
	streams string
	entries string
	showChaps bool
	format string
}

func NewArgs() *Args {
	return &Args{
		Flags: Flags{Profile: "default"},
	}
}

type cmdArgs struct {
	args []string
}

func(arg *cmdArgs) Append(args ...string) {
	arg.args = append(arg.args, args...)
}

func(cmd Cmd) ParseArgs() *Cmd {
	if log := cmd.args.LogLevel; log != "" {
		cmd.appendArgs("-v", log)
	}

	// parse cmd args
	switch {
	case cmd.ffmpeg:
		cmd.parseFFmpegArgs()
	case cmd.ffprobe:
		cmd.parseFFprobeArgs()
	}

	switch {
	case cmd.ffmpeg:
		cmd.exec = exec.Command("ffmpeg", cmd.cmdArgs...)
	case cmd.ffprobe:
		//cmd.exec = exec.Command("ffprobe", cmd.cmdArgs...)
	}
	return &cmd
}

func(cmd *Cmd) appendArgs(args ...string) *Cmd {
	cmd.cmdArgs = append(cmd.cmdArgs, args...)
	return cmd
}

func(cmd *Cmd) parseFFmpegArgs() *Cmd {
	if cmd.args.Overwrite {
		cmd.appendArgs("-y")
	}

	// pre input
	if pre := cmd.args.PreInput; len(pre) > 0 {
		cmd.appendArgs(pre.Split()...)
	}

	// input
	cmd.appendArgs("-i", cmd.Input)
	cover := cmd.Flags.CoverFile
	meta := cmd.Flags.MetaFile
	if meta != "" {
		cmd.appendArgs("-i", meta)
	}

	if cover != "" {
		cmd.appendArgs("-i", cover)
	}

	//map input
	idx := 0
	if cover != "" || meta != "" {
		cmd.appendArgs("-map", strconv.Itoa(idx) + ":0")
		idx++
	}

	if cover != "" {
		cmd.appendArgs("-map", "0:" + strconv.Itoa(idx))
		idx++
	}

	if meta != "" {
		cmd.appendArgs("-map_metadata", strconv.Itoa(idx))
		idx++
	}

	// post input
	if post := cmd.args.PostInput; len(post) > 0 {
		cmd.appendArgs(post.Split()...)
	}

	//video codec
	if codec := cmd.args.VideoCodec; codec != "" {
		switch codec {
		case "":
		case "none", "vn":
			cmd.appendArgs("-vn")
		default:
			cmd.appendArgs("-c:v", codec)
			//video params
			if params := cmd.args.VideoParams.Split(); len(params) > 0 {
				cmd.appendArgs(params...)
			}

			//video filters
			if filters := cmd.args.VideoFilters.Join(); len(filters) > 0 {
				cmd.appendArgs("-vf", filters)
			}
		}
	}

	//filter complex
	if filters := cmd.args.FilterComplex.Join(); len(filters) > 0 {
		cmd.appendArgs("-vf", filters)
	}

	//audio codec
	if codec := cmd.args.AudioCodec; codec != "" {
		switch codec {
		case "":
		case "none", "an":
			cmd.appendArgs("-an")
		default:
			cmd.appendArgs("-c:a", codec)
			//audio params
			if params := cmd.args.AudioParams.Split(); len(params) > 0 {
				cmd.appendArgs(params...)
			}

			//audio filters
			if filters := cmd.args.AudioFilters.Join(); len(filters) > 0 {
				cmd.appendArgs("-af", filters)
			}
		}
	}

	//output
	var (
		name string
		ext string
	)

	if out := cmd.args.Output; out != "" {
		name = out
	}

	if p := cmd.args.Padding; p != "" {
		name = name + fmt.Sprintf(p, cmd.args.num)
	}

	switch {
	case cmd.Ext != "":
		ext = cmd.Ext
	case cmd.args.Ext != "":
		ext = cmd.args.Ext
	default:
		ext = cmd.Media.Ext
	}
	cmd.appendArgs(name + ext)

	return cmd
}

func(cmd *Cmd) parseFFprobeArgs() *Cmd {
	if log := cmd.args.LogLevel; log != "" {
		cmd.appendArgs("-v", log)
	}

	if cmd.args.pretty {
		cmd.appendArgs("-pretty")
	}

	if stream := cmd.args.streams; stream != "" {
		cmd.appendArgs("-select_streams", stream)
	}

	if entries := cmd.args.entries; entries != "" {
		cmd.appendArgs("-show_entries", entries)
	}

	if cmd.args.showChaps {
		cmd.appendArgs("-show_chapters")
	}

	cmd.appendArgs("-of")
	switch f := cmd.args.format; f {
	default:
		fallthrough
	case "":
		fallthrough
	case "plain":
		cmd.appendArgs("default=noprint_wrappers=1:nokey=1")
	case "json":
		cmd.appendArgs("json=c=1")
	}

	cmd.appendArgs(cmd.Input)

	//cmd.exec = exec.Command("ffprobe", cmd.cmdArgs...)
	return cmd
}

