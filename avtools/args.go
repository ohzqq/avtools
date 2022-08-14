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
	cmdArgs
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
	CoverFile   string
	MetaFile    string
	CueFile     string
}

func NewArgs() *Args {
	return &Args{
		Options: Options{Profile: "default"},
	}
}

//func (args *Args) Output() *Args {
//  return args
//}

func (args *Args) parseInput() *Args {
	if args.Input != "" {
		args.Append("-i", args.Input)
	}

	meta := args.MetaFile
	if meta != "" {
		args.Append("-i", meta)
	}

	cover := args.CoverFile
	if cover != "" {
		args.Append("-i", cover)
	}

	//map input
	idx := 0
	if cover != "" || meta != "" {
		args.Append("-map", strconv.Itoa(idx)+":0")
		idx++
	}

	if cover != "" {
		args.Append("-map", strconv.Itoa(idx)+":0")
		idx++
	}

	if meta != "" {
		args.Append("-map_metadata", strconv.Itoa(idx))
		idx++
	}

	return args
}

func (args *Args) videoArgs() *Args {
	//video codec
	if codec := args.VideoCodec; codec != "" {
		switch codec {
		case "":
		case "none", "vn":
			args.Append("-vn")
		default:
			args.Append("-c:v", codec)
			//video params
			if params := args.VideoParams.Split(); len(params) > 0 {
				args.Append(params...)
			}

			//video filters
			if filters := args.VideoFilters.Join(); len(filters) > 0 {
				args.Append("-vf", filters)
			}
		}
	}

	return args
}

func (args *Args) audioArgs() *Args {
	//audio codec
	if codec := args.AudioCodec; codec != "" {
		switch codec {
		case "":
		case "none", "an":
			args.Append("-an")
		default:
			args.Append("-c:a", codec)
			//audio params
			if params := args.AudioParams.Split(); len(params) > 0 {
				args.Append(params...)
			}

			//audio filters
			if filters := args.AudioFilters.Join(); len(filters) > 0 {
				args.Append("-af", filters)
			}
		}
	}

	return args
}

func (args *Args) output() *Args {
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

	args.Append(name + ext)

	return args
}

func (args *Args) Parse() []string {
	if log := args.LogLevel; log != "" {
		args.Append("-v", log)
	}

	if args.Overwrite {
		args.Append("-y")
	}

	args.parseInput()

	// pre input
	if pre := args.PreInput; len(pre) > 0 {
		args.Append(pre.Split()...)
	}

	// post input
	if post := args.PostInput; len(post) > 0 {
		args.Append(post.Split()...)
	}

	//filter complex
	if filters := args.FilterComplex.Join(); len(filters) > 0 {
		args.Append("-vf", filters)
	}

	args.videoArgs()
	args.audioArgs()
	args.output()

	return args.cmdArgs.args
}
