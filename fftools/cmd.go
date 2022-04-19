package fftools

import (
	//"os"
	"strings"
	"os/exec"

	//"github.com/alessio/shellescape"
)

type CmdArgs map[string]string

var cmdFlags = map[string]string{
	//"Pre": "pre",
	//"Input": "i",
	//"Post": "post",
	//"VideoCodec": "c:v",
	//"VideoParams": "vp",
	//"VideoFilters": "vf",
	//"AudioCodec": "c:a",
	//"AudioParams": "ap",
	//"AudioFilters": "af",
	//"FilterCompex": "filter",
	"Output": "o",
	"Cue": "cue",
	"Cover": "c",
	"FFmetadata": "m",
	//"Verbosity": "v",
	"Profile": "p",
}

var argOrder = []string{"Verbosity", "Overwrite", "Pre", "Input", "Post", "VideoCodec", "VideoParams", "VideoFilters", "FilterComplex", "AudioCodec", "AudioParams", "AudioFilters", "Output", "Ext"}

type FFmpegCmd struct {
	cmd *exec.Cmd
	input []string
	Arguments CmdArgs
}

func NewCmd() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.cmd = exec.Command("ffmpeg", "-hide_banner")
	ff.Arguments = make(CmdArgs)
	return &ff
}

func (ff *FFmpegCmd) Args(args CmdArgs) *FFmpegCmd {
	ff.Arguments = args
	//ff.Arguments.SetVerbosity()
	return ff
}

func (ff *FFmpegCmd) Input(input []string) *FFmpegCmd {
	var in []string
	for _, i := range input {
		in = append(in, "-i ")
		in = append(in, i)
	}
	ff.Arguments["Input"] = strings.Join(in, " ")
	return ff
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Verbosity())
		case "Overwrite":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Overwrite())
		case "Pre":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Pre())
		case "Input":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Input())
		case "Post":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Post())
		case "VideoCodec":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.VC())
		case "VideoParams":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.VP())
		case "VideoFilters":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.VF())
		case "FilterComplex":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Filter())
		case "AudioCodec":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.AC())
		case "AudioParams":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.AP())
		case "AudioFilters":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.AF())
		case "Output":
			ff.cmd.Args = append(ff.cmd.Args, ff.Arguments.Output())
		}
	}
	return ff.cmd
}

func (c CmdArgs) Verbosity() string {
	var v string
	if Cfg.Defaults.Verbosity != "" {
		v = "-loglevel " + Cfg.Defaults.Verbosity
	}
	return v
}

func (c CmdArgs) Pre() string {
	return c["Pre"]
}

func (c CmdArgs) Overwrite() string {
	var o string
	if Cfg.Defaults.Overwrite {
			o = "-y"
	} else {
			o = ""
	}
	return o
}

func (c CmdArgs) Input() string {
	return c["Input"]
}

func (c CmdArgs) Post() string {
	return c["Post"]
}

func (c CmdArgs) VC() string {
	if c["VideoCodec"] != "" {
		return "-c:v " + c["VideoCodec"]
	}
	return ""
}

func (c CmdArgs) VP() string {
	return c["VideoParams"]
}

func (c CmdArgs) VF() string {
	if c["VideoFilters"] != "" {
		return "-vf " + c["VideoFilters"]
	}
	return ""
}

func (c CmdArgs) Filter() string {
	if c["FilterComplex"] != "" {
		return "-filter " + c["FilterComplex"]
	}
	return ""
}

func (c CmdArgs) AC() string {
	if c["AudioCodec"] != "" {
		return "-c:a " + c["AudioCodec"]
	}
	return ""
}

func (c CmdArgs) AP() string {
	return c["AudioParams"]
}

func (c CmdArgs) AF() string {
	if c["AudioFilters"] != "" {
		return "-af " + c["AudioFilters"]
	}
	return ""
}

func (c CmdArgs) Output() string {
	var o string
	var pad string
	var ext string
	if Cfg.Defaults.Output != "" {
		o = Cfg.Defaults.Output
		if Cfg.Defaults.Padding {
			pad = "%06d"
		} else {
			pad = ""
		}
	}
	if c["Ext"] != "" {
		ext = "." + c["Ext"]
	} else {
		ext = ".mkv"
	}
	return o + pad + ext
}
