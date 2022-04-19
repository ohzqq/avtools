package fftools

import (
	//"os"
	"strings"
	"strconv"
	"os/exec"

	//"github.com/alessio/shellescape"
)

type FFmpegCmd struct {
	cmd *exec.Cmd
	Input []string
	profile bool
	args CmdArgs
}

func NewCmd() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) Profile(p string) *FFmpegCmd {
	ff.profile = true
	ff.args = Cfg.Profiles[p]
	ff.args.LogLevel(Cfg.Defaults.Verbosity)
	ff.args.Out(Cfg.Defaults.Output)
	ff.args.OverWrite(Cfg.Defaults.Overwrite)
	return ff
}

func (ff *FFmpegCmd) Args() *CmdArgs {
	if !ff.profile {
		ff.args = CmdArgs{}
		ff.args.VCodec("copy")
		ff.args.ACodec("copy")
	}
	return &ff.args
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	argOrder := []string{"Verbosity", "Overwrite", "Pre", "Input", "Meta", "Post", "VideoCodec", "VideoParams", "VideoFilters", "FilterComplex", "AudioCodec", "AudioParams", "AudioFilters", "Output", "Ext"}
	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			if ff.Verbosity() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Verbosity())
			}
		case "Overwrite":
			if ff.Overwrite() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Overwrite())
			}
		case "Pre":
			if ff.Pre() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Pre())
			}
		case "Input":
			if len(ff.Input) > 0 {
				ff.cmd.Args = append(ff.cmd.Args, ff.joinInput())
			}
		case "Meta":
			if ff.Meta() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Meta())
			}
		case "Post":
			if ff.Post() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Post())
			}
		case "VideoCodec":
			if ff.VideoCodec() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.VideoCodec())
			}
		case "VideoParams":
			if ff.VideoParams() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.VideoParams())
			}
		case "VideoFilters":
			if ff.VideoFilters() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.VideoFilters())
			}
		case "FilterComplex":
			if ff.FilterComplex() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.FilterComplex())
			}
		case "AudioCodec":
			if ff.AudioCodec() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.AudioCodec())
			}
		case "AudioParams":
			if ff.AudioParams() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.AudioParams())
			}
		case "AudioFilters":
			if ff.AudioFilters() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.AudioFilters())
			}
		case "Output":
			if ff.Output() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Output())
			}
		}
	}
	return ff.cmd
}

func (ff *FFmpegCmd) Verbosity() string {
	var v string
	if Cfg.Defaults.Verbosity == "" {
		v = ""
	} else {
		v = "-loglevel " + Cfg.Defaults.Verbosity
	}
	return v
}

func (ff *FFmpegCmd) In(input string) {
	ff.Input = append(ff.Input, input)
}

func (ff *FFmpegCmd) joinInput() string {
	var in []string
	for _, i := range ff.Input {
		in = append(in, "-i")
		in = append(in, i)
	}
	return strings.Join(in, " ")
}


func (ff *FFmpegCmd) Meta() string {
	idx := strconv.Itoa(len(ff.Input))
	if ff.args.Metadata != "" {
		return "-i " + ff.args.Metadata + " -map_metadata " + idx
	}
	return ""
}

func (ff *FFmpegCmd) Pre() string {
	return ff.args.Pre
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

func (ff *FFmpegCmd) Post() string {
	return ff.args.Post
}

func (ff *FFmpegCmd) VideoCodec() string {
	if ff.args.VideoCodec != "" {
		return "-c:v " + ff.args.VideoCodec
	}
	return ""
}

func (ff *FFmpegCmd) VideoParams() string {
	return ff.args.VideoParams
}

func (ff *FFmpegCmd) VideoFilters() string {
	if ff.args.VideoFilters != "" {
		return "-vf " + ff.args.VideoFilters
	}
	return ""
}

func (ff *FFmpegCmd) FilterComplex() string {
	if ff.args.FilterComplex != "" {
		return "-filter " + ff.args.FilterComplex
	}
	return ""
}

func (ff *FFmpegCmd) AudioCodec() string {
	if ff.args.AudioCodec != "" {
		return "-c:a " + ff.args.AudioCodec
	}
	return ""
}

func (ff *FFmpegCmd) AudioParams() string {
	return ff.args.AudioParams
}

func (ff *FFmpegCmd) AudioFilters() string {
	if ff.args.AudioFilters != "" {
		return "-af " + ff.args.AudioFilters
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
	if ff.args.Ext != "" {
		ext = "." + ff.args.Ext
	} else {
		ext = ".mkv"
	}
	if ff.args.Output == "" {
		return o + pad + ext
	} else {
		return ff.args.Output + pad + ext
	}
	return ""
}
