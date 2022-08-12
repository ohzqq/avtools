package avtools

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Args struct {
	Options
	Input         string
	PreInput      mapArgs
	PostInput     mapArgs
	VideoCodec    string
	VideoParams   mapArgs
	VideoFilters  stringArgs
	AudioCodec    string
	AudioParams   mapArgs
	AudioFilters  stringArgs
	FilterComplex stringArgs
	MiscParams    stringArgs
	LogLevel      string
	Name          string
	Padding       string
	Ext           string
	num           int
	pretty        bool
	streams       string
	entries       string
	showChaps     bool
	format        string
}

type Options struct {
	Overwrite   bool
	Profile     string
	Start       string
	End         string
	Output      string
	ChapNo      int
	MetaSwitch  bool
	CoverSwitch bool
	CueSwitch   bool
	ChapSwitch  bool
	Verbose     bool
	CoverFile   string
	MetaFile    string
	CueFile     string
}

func NewArgs() *Args {
	return &Args{
		Options: Options{Profile: "default"},
	}
}

func (cmd *ffmpegCmd) ParseOptions() *ffmpegCmd {
	cmd.Args = Cfg().GetProfile(cmd.opts.Profile)

	if meta := cmd.opts.MetaFile; meta != "" {
		cmd.media.SetFile("ffmeta", meta)
		//cmd.media.AddFFmeta(meta)
	}

	if cover := cmd.opts.CoverFile; cover != "" {
		cmd.media.SetFile("cover", cover)
		//cmd.media.AddCover(cover)
	}

	if cue := cmd.opts.CueFile; cue != "" {
		cmd.media.SetFile("cue", cue)
		//cmd.media.AddCue(cue)
	}

	if y := cmd.opts.Overwrite; y {
		cmd.Overwrite = y
	}

	if o := cmd.opts.Output; o != "" {
		cmd.Name = o
	}

	if c := cmd.opts.ChapNo; c != 0 {
		cmd.num = c
	}

	return cmd
}

func (cmd *ffmpegCmd) ParseArgs() *Cmd {
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

	m := cmd.media.GetFile("input")
	if cmd.media != nil {
		cmd.args.Append("-i", m.Path)
	}

	if cmd.Input != "" {
		cmd.args.Append("-i", cmd.Input)
	}

	meta := cmd.opts.MetaFile
	if meta != "" {
		cmd.args.Append("-i", meta)
	}

	cover := cmd.opts.CoverFile
	if cover != "" {
		cmd.args.Append("-i", cover)
	}

	//map input
	idx := 0
	if cover != "" || meta != "" {
		cmd.args.Append("-map", strconv.Itoa(idx)+":0")
		idx++
	}

	if cover != "" {
		cmd.args.Append("-map", strconv.Itoa(idx)+":0")
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
		ext  string
	)

	if out := cmd.Output; out != "" {
		name = out
	}

	if p := cmd.Padding; p != "" {
		name = name + fmt.Sprintf(p, cmd.num)
	}

	media := cmd.media.GetFile("input")
	switch {
	case cmd.Ext != "":
		ext = cmd.Ext
	default:
		ext = media.Ext
	}
	cmd.args.Append(name + ext)

	return NewCmd(exec.Command("ffmpeg", cmd.args.args...), cmd.opts.Verbose)
}

type cmdArgs struct {
	args []string
}

func (arg *cmdArgs) Append(args ...string) {
	arg.args = append(arg.args, args...)
}

type mapArgs []map[string]string

func newMapArg(k, v string) map[string]string {
	return map[string]string{k: v}
}

func (a *Args) AppendMapArg(key, flag, value string) {
	mapArg := newMapArg(flag, value)
	switch key {
	case "pre":
		a.PreInput = append(a.PreInput, mapArg)
	case "post":
		a.PostInput = append(a.PostInput, mapArg)
	case "videoParams":
		a.VideoParams = append(a.VideoParams, mapArg)
	case "audioParams":
		a.AudioParams = append(a.AudioParams, mapArg)
	}
}

func (m mapArgs) Split() []string {
	var args []string
	for _, flArg := range m {
		for flag, arg := range flArg {
			args = append(args, "-"+flag, arg)
		}
	}
	return args
}

type stringArgs []string

func (s stringArgs) Join() string {
	return strings.Join(s, ",")
}
