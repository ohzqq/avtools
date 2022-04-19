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
	FilterComplex string
	Verbosity string
	Output string
	Padding bool
	Ext string
	Overwrite bool
}

var argOrder = []string{"Verbosity", "Overwrite", "Pre", "Input", "Post", "VideoCodec", "VideoParams", "VideoFilters", "FilterComplex", "AudioCodec", "AudioParams", "AudioFilters", "Output", "Ext"}

type FFmpegCmd struct {
	cmd *exec.Cmd
	input []string
	Args CmdArgs
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
	//ff.Args.Input = strings.Join(in, " ")
	ff.Args.Input = strings.Join(in, " ")
	return ff
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.cmd.Args = append(ff.cmd.Args, ff.Verbosity())
		case "Overwrite":
			ff.cmd.Args = append(ff.cmd.Args, ff.Overwrite())
		case "Pre":
			ff.cmd.Args = append(ff.cmd.Args, ff.Pre())
		case "Input":
			//ff.cmd.Args = append(ff.cmd.Args, ff.Input())
		case "Post":
			ff.cmd.Args = append(ff.cmd.Args, ff.Post())
		case "VideoCodec":
			ff.cmd.Args = append(ff.cmd.Args, ff.VC())
		case "VideoParams":
			ff.cmd.Args = append(ff.cmd.Args, ff.VP())
		case "VideoFilters":
			ff.cmd.Args = append(ff.cmd.Args, ff.VF())
		case "FilterComplex":
			ff.cmd.Args = append(ff.cmd.Args, ff.Filter())
		case "AudioCodec":
			ff.cmd.Args = append(ff.cmd.Args, ff.AC())
		case "AudioParams":
			ff.cmd.Args = append(ff.cmd.Args, ff.AP())
		case "AudioFilters":
			ff.cmd.Args = append(ff.cmd.Args, ff.AF())
		case "Output":
			ff.cmd.Args = append(ff.cmd.Args, ff.Output())
		}
	}
	return ff.cmd
}

func (ff *FFmpegCmd) Verbosity() string {
	var v string
	if Cfg.Defaults.Verbosity != "" {
		v = "-loglevel " + Cfg.Defaults.Verbosity
	}
	return v
}

func (ff *FFmpegCmd) Pre() string {
	return ff.Args.Pre
}

func (ff *FFmpegCmd) Overwrite() string {
	var o string
	if Cfg.Defaults.Overwrite {
			o = "-y"
	} else {
			o = ""
	}
	return o
}

//func (ff *FFmpegCmd) Input() string {
//  return ff.Args.Input
//}

func (ff *FFmpegCmd) Post() string {
	return ff.Args.Post
}

func (ff *FFmpegCmd) VC() string {
	if ff.Args.VideoCodec != "" {
		return "-c:v " + ff.Args.VideoCodec
	}
	return ""
}

func (ff *FFmpegCmd) VP() string {
	return ff.Args.VideoParams
}

func (ff *FFmpegCmd) VF() string {
	if ff.Args.VideoFilters != "" {
		return "-vf " + ff.Args.VideoFilters
	}
	return ""
}

func (ff *FFmpegCmd) Filter() string {
	if ff.Args.FilterComplex != "" {
		return "-filter " + ff.Args.FilterComplex
	}
	return ""
}

func (ff *FFmpegCmd) AC() string {
	if ff.Args.AudioCodec != "" {
		return "-c:a " + ff.Args.AudioCodec
	}
	return ""
}

func (ff *FFmpegCmd) AP() string {
	return ff.Args.AudioParams
}

func (ff *FFmpegCmd) AF() string {
	if ff.Args.AudioFilters != "" {
		return "-af " + ff.Args.AudioFilters
	}
	return ""
}

func (ff *FFmpegCmd) Output() string {
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
	if ff.Args.Ext != "" {
		ext = "." + ff.Args.Ext
	} else {
		ext = ".mkv"
	}
	return o + pad + ext
}
