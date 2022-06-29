package avtools

import (
	"fmt"
	//"os"
	"os/exec"
	//"bytes"
	//"log"
	"strconv"
	"strings"
	//"path/filepath"
)
var _ = fmt.Printf

type ffmpegArgs struct {
	*Flags
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
}

type mapArgs []map[string]string

func newMapArg(k, v string) map[string]string {
	return map[string]string{k: v}
}

func(m mapArgs) Split() []string {
	var args []string
	for _, flArg := range m {
		for flag, arg := range flArg {
			args = append(args, "-" + flag, arg)
		}
	}
	return args
}

type stringArgs []string

func(s stringArgs) Join() string {
	return strings.Join(s, ",")
}

type ffmpegCmd struct {
	Input string
	media *Media
	args cmdArgs
	exec *exec.Cmd
	*ffmpegArgs
}

func NewFFmpegCmd(i string) *ffmpegCmd {
	return &ffmpegCmd{
		Input: i,
		media: NewMedia(i),
	}
}

func(cmd *ffmpegCmd) SetFlags(f *Flags) *ffmpegCmd {
	fmt.Printf("%+v\n", f)
	cmd.ffmpegArgs.Flags = f
	return cmd
}

func(cmd *ffmpegCmd) String() string {
	return cmd.exec.String()
}

func(cmd *ffmpegCmd) ParseFlags() *ffmpegCmd {
	cmd.ffmpegArgs = Cfg().GetProfile(cmd.Profile)
	cmd.media.JsonMeta().Unmarshal()

	if meta := cmd.MetaFile; meta != "" {
		cmd.media.SetMeta(LoadFFmetadataIni(meta))
	}

	if cue := cmd.CueFile; cue != "" {
		cmd.media.SetChapters(LoadCueSheet(cue))
	}

	if y := cmd.Overwrite; y {
		cmd.Overwrite = y
	}

	if o := cmd.Output; o != "" {
		cmd.Name = o
	}

	if c := cmd.ChapNo; c  != 0 {
		cmd.num = c
	}

	return cmd
}

func(cmd *ffmpegCmd) Parse() Cmd {
	if log := cmd.LogLevel; log != "" {
		cmd.args.Append("-v", log)
	}

	if cmd.Overwrite {
		cmd.args.Append("-y")
	}

	// pre input
	if pre := cmd.PreInput; len(pre) > 0 {
		cmd.args.Append(pre.Split()...)
	}

	// input
	cmd.args.Append("-i", cmd.Input)
	cover := cmd.Flags.CoverFile
	meta := cmd.Flags.MetaFile
	if meta != "" {
		cmd.args.Append("-i", meta)
	}

	if cover != "" {
		cmd.args.Append("-i", cover)
	}

	//map input
	idx := 0
	if cover != "" || meta != "" {
		cmd.args.Append("-map", strconv.Itoa(idx) + ":0")
		idx++
	}

	if cover != "" {
		cmd.args.Append("-map", "0:" + strconv.Itoa(idx))
		idx++
	}

	if meta != "" {
		cmd.args.Append("-map_metadata", strconv.Itoa(idx))
		idx++
	}

	// post input
	if post := cmd.PostInput; len(post) > 0 {
		cmd.args.Append(post.Split()...)
	}

	//video codec
	if codec := cmd.VideoCodec; codec != "" {
		switch codec {
		case "":
		case "none", "vn":
			cmd.args.Append("-vn")
		default:
			cmd.args.Append("-c:v", codec)
			//video params
			if params := cmd.VideoParams.Split(); len(params) > 0 {
				cmd.args.Append(params...)
			}

			//video filters
			if filters := cmd.VideoFilters.Join(); len(filters) > 0 {
				cmd.args.Append("-vf", filters)
			}
		}
	}

	//filter complex
	if filters := cmd.FilterComplex.Join(); len(filters) > 0 {
		cmd.args.Append("-vf", filters)
	}

	//audio codec
	if codec := cmd.AudioCodec; codec != "" {
		switch codec {
		case "":
		case "none", "an":
			cmd.args.Append("-an")
		default:
			cmd.args.Append("-c:a", codec)
			//audio params
			if params := cmd.AudioParams.Split(); len(params) > 0 {
				cmd.args.Append(params...)
			}

			//audio filters
			if filters := cmd.AudioFilters.Join(); len(filters) > 0 {
				cmd.args.Append("-af", filters)
			}
		}
	}

	//output
	var (
		name string
		ext string
	)

	if out := cmd.Output; out != "" {
		name = out
	}

	if p := cmd.Padding; p != "" {
		name = name + fmt.Sprintf(p, cmd.num)
	}

	switch {
	//case cmd.Ext != "":
	//  ext = cmd.Ext
	case cmd.Ext != "":
		ext = cmd.Ext
	default:
		ext = cmd.media.Ext
	}
	cmd.args.Append(name + ext)

	return Cmd{exec: exec.Command("ffmpeg", cmd.args.args...)}
}
