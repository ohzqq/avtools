package avtools

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

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

type Args struct {
	Options
	input         []string
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
	JsonSwitch  bool
	Verbose     bool
	Input       string
	CoverFile   string
	MetaFile    string
	CueFile     string
}

func NewArgs() *Args {
	return &Args{}
}

func (args *Args) Parse(opts Options) []string {
	args.Options = opts

	cmdArgs := cmdArgs{}
	if log := args.LogLevel; log != "" {
		cmdArgs.Append("-v", log)
	}

	if opts.Overwrite {
		cmdArgs.Append("-y")
	}

	for _, i := range opts.parseInput() {
		cmdArgs.Append(i)
	}

	// pre input
	if pre := args.PreInput; len(pre) > 0 {
		cmdArgs.Append(pre.Split()...)
	}

	// post input
	if post := args.PostInput; len(post) > 0 {
		cmdArgs.Append(post.Split()...)
	}

	//filter complex
	if filters := args.FilterComplex.Join(); len(filters) > 0 {
		cmdArgs.Append("-vf", filters)
	}

	args.videoArgs(&cmdArgs)
	args.audioArgs(&cmdArgs)
	args.output(&cmdArgs)

	return cmdArgs.args
}

//func (args *Args) Output() *Args {
//  return args
//}

func (opts Options) parseInput() []string {
	var input []string

	if opts.Input != "" {
		input = append(input, "-i", opts.Input)
	}

	meta := opts.MetaFile
	if meta != "" {
		input = append(input, "-i", meta)
	}

	cover := opts.CoverFile
	if cover != "" {
		input = append(input, "-i", cover)
	}

	idx := 0
	if cover != "" || meta != "" {
		input = append(input, "-map", strconv.Itoa(idx)+":0")
		idx++
	}

	if cover != "" {
		input = append(input, "-map", strconv.Itoa(idx)+":0")
		idx++
	}

	if meta != "" {
		input = append(input, "-map_metadata", strconv.Itoa(idx))
		idx++
	}

	return input
}

func (args *Args) videoArgs(a *cmdArgs) {
	//video codec
	if codec := args.VideoCodec; codec != "" {
		switch codec {
		case "":
		case "none", "vn":
			a.Append("-vn")
		default:
			a.Append("-c:v", codec)
			//video params
			if params := args.VideoParams.Split(); len(params) > 0 {
				a.Append(params...)
			}

			//video filters
			if filters := args.VideoFilters.Join(); len(filters) > 0 {
				a.Append("-vf", filters)
			}
		}
	}
}

func (args *Args) audioArgs(a *cmdArgs) {
	//audio codec
	if codec := args.AudioCodec; codec != "" {
		switch codec {
		case "":
		case "none", "an":
			a.Append("-an")
		default:
			a.Append("-c:a", codec)
			//audio params
			if params := args.AudioParams.Split(); len(params) > 0 {
				a.Append(params...)
			}

			//audio filters
			if filters := args.AudioFilters.Join(); len(filters) > 0 {
				a.Append("-af", filters)
			}
		}
	}
}

func (args *Args) output(a *cmdArgs) {
	//output
	var (
		name = args.Name
		ext  = filepath.Ext(args.Input)
	)

	if out := args.Output; out != "" {
		name = out
	}

	if p := args.Padding; p != "" {
		name = name + fmt.Sprintf(p, args.num)
	}

	if e := args.Ext; e != "" {
		ext = e
	}

	a.Append(name + ext)

	//return a
}
