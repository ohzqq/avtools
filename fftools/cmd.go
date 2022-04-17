package fftools

import (
	//"os"
	"strings"
	"os/exec"
)

type FFmpegCmd struct {
	Cmd *exec.Cmd
	Input string
	VideoFilter string
	AudioFilter string
	CodecVideo string
	CodecAudio string
	VideoParams string
	AudioParams string
}

func New() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.Cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) Args(args ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, args...)
}

func (ff *FFmpegCmd) Files(input ...string) string {
	return strings.Join(input, " ")
}

func (ff *FFmpegCmd) VF(filters ...string) string {
	return strings.Join(filters , " ")
}

func (ff *FFmpegCmd) AF(filters ...string) string {
	return strings.Join(filters , " ")
}

func (ff *FFmpegCmd) FilterComplex(filter string) string {
	return filter
}

func (ff *FFmpegCmd) CV(codec string) string {
	return codec
}

func (ff *FFmpegCmd) CA(codec string) string {
	return codec
}

func (ff *FFmpegCmd) VParams(params ...string) string {
	return strings.Join(params, " ")
}

func (ff *FFmpegCmd) AParams(params ...string) string {
	return strings.Join(params, " ")
}

func (p profile) Start(s string) profile {
	var in []string
	for _, arg := range p {
		if arg == "ss" {
			in = append(in, arg + " " + s)
		} else {
			in = append(in, arg)
		}
	}
	return profile(in)
}

func (p profile) End(e string) profile {
	var in []string
	for _, arg := range p {
		if arg == "to" || arg == "t" {
			in = append(in, arg + " " + e)
		} else {
			in = append(in, arg)
		}
	}
	return profile(in)
}

func (p profile) Input(input ...string) profile {
	var in []string
	i := 0
	for _, arg := range p {
		if arg == "i" {
			in = append(in, arg + " " + input[i])
			i++
		} else {
			in = append(in, arg)
		}
	}
	return profile(in)
}


func (p profile) String() string {
	var cmdString strings.Builder
	for _, arg := range p {
		cmdString.WriteString("-")
		cmdString.WriteString(arg)
		cmdString.WriteString(" ")
	}
	return cmdString.String()
}
