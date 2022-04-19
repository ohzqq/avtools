package fftools

import (
	//"os"
	"strings"
	"os/exec"

	//"github.com/alessio/shellescape"
)

//type CmdArgs map[string]string
type CmdArgs struct {
	Pre string
	Input string
	Post string
	VideoCodec string
	VideoParams string
	VideoFilters string
	AudioCodec string
	AudioParams string
	AudioFilters string
	FilterCompex string
	Verbosity string
	Output string
	Padding bool
	Ext string
	Overwrite bool
}


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
	CmdArgs
}

func NewCmd() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.cmd = exec.Command("ffmpeg", "-hide_banner")
	ff.Args = CmdArgs{}
	return &ff
}

//func (ff *FFmpegCmd) Args(args CmdArgs) *FFmpegCmd {
//  ff.Args = args
//  //ff.Args.SetVerbosity()
//  return ff
//}

func (ff *FFmpegCmd) Input(input []string) *FFmpegCmd {
	var in []string
	for _, i := range input {
		in = append(in, "-i ")
		in = append(in, i)
	}
	//ff.Args["Input"] = strings.Join(in, " ")
	ff.Args["Input"] = strings.Join(in, " ")
	return ff
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Verbosity())
		case "Overwrite":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Overwrite())
		case "Pre":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Pre())
		case "Input":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Input())
		case "Post":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Post())
		case "VideoCodec":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.VC())
		case "VideoParams":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.VP())
		case "VideoFilters":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.VF())
		case "FilterComplex":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Filter())
		case "AudioCodec":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.AC())
		case "AudioParams":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.AP())
		case "AudioFilters":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.AF())
		case "Output":
			ff.cmd.Args = append(ff.cmd.Args, ff.Args.Output())
		}
	}
	return ff.cmd
}

func (a *FFmpegCmd) Verbosity() string {
	var v string
	if Cfg.Defaults.Verbosity != "" {
		v = "-loglevel " + Cfg.Defaults.Verbosity
	}
	return v
}

func (a *FFmpegCmd) Pre() string {
	return c["Pre"]
}

func (a *FFmpegCmd) Overwrite() string {
	var o string
	if Cfg.Defaults.Overwrite {
			o = "-y"
	} else {
			o = ""
	}
	return o
}

func (a *FFmpegCmd) Input() string {
	return c["Input"]
}

func (a *FFmpegCmd) Post() string {
	return c["Post"]
}

func (a *FFmpegCmd) VC() string {
	if c["VideoCodec"] != "" {
		return "-c:v " + c["VideoCodec"]
	}
	return ""
}

func (a *FFmpegCmd) VP() string {
	return c["VideoParams"]
}

func (a *FFmpegCmd) VF() string {
	if c["VideoFilters"] != "" {
		return "-vf " + c["VideoFilters"]
	}
	return ""
}

func (a *FFmpegCmd) Filter() string {
	if c["FilterComplex"] != "" {
		return "-filter " + c["FilterComplex"]
	}
	return ""
}

func (a *FFmpegCmd) AC() string {
	if c["AudioCodec"] != "" {
		return "-c:a " + c["AudioCodec"]
	}
	return ""
}

func (a *FFmpegCmd) AP() string {
	return c["AudioParams"]
}

func (a *FFmpegCmd) AF() string {
	if c["AudioFilters"] != "" {
		return "-af " + c["AudioFilters"]
	}
	return ""
}

func (a *FFmpegCmd) Output() string {
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
