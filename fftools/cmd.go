package fftools

import (
	//"os"
	//"strings"
	"os/exec"
)

var ArgOrder = []string{"Pre", "Input", "Post", "VideoCodec", "VideoParams", "VideoFilters", "AudioCodec", "AudioParams", "AudioFilters", "FilterComplex", "Output"}

type FFmpegCmd struct {
	Cmd *exec.Cmd
	Input string
	Profile Profile
}

func New() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.Cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) PreInput(args ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, args...)
}

func (ff *FFmpegCmd) Files(input ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, input...)
}

func (ff *FFmpegCmd) PostInput(args ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, args...)
}

func (ff *FFmpegCmd) CV(codec string) {
	ff.Cmd.Args = append(ff.Cmd.Args, codec)
}

func (ff *FFmpegCmd) VParams(params ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, params...)
}

func (ff *FFmpegCmd) VF(filters ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, filters...)
}

func (ff *FFmpegCmd) FilterComplex(filter string) {
	ff.Cmd.Args = append(ff.Cmd.Args, filter)
}

func (ff *FFmpegCmd) CA(codec string) {
	ff.Cmd.Args = append(ff.Cmd.Args, codec)
}

func (ff *FFmpegCmd) AParams(params ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, params...)
}

func (ff *FFmpegCmd) AF(filters ...string) {
	ff.Cmd.Args = append(ff.Cmd.Args, filters...)
}

//func (p profile) Start(s string) profile {
//  var in []string
//  for _, arg := range p {
//    if arg == "ss" {
//      in = append(in, arg + " " + s)
//    } else {
//      in = append(in, arg)
//    }
//  }
//  return profile(in)
//}

//func (p profile) End(e string) profile {
//  var in []string
//  for _, arg := range p {
//    if arg == "to" || arg == "t" {
//      in = append(in, arg + " " + e)
//    } else {
//      in = append(in, arg)
//    }
//  }
//  return profile(in)
//}

//func (p profile) Input(input ...string) profile {
//  var in []string
//  i := 0
//  for _, arg := range p {
//    if arg == "i" {
//      in = append(in, arg + " " + input[i])
//      i++
//    } else {
//      in = append(in, arg)
//    }
//  }
//  return profile(in)
//}


//func (p profile) String() string {
//  var cmdString strings.Builder
//  for _, arg := range p {
//    cmdString.WriteString("-")
//    cmdString.WriteString(arg)
//    cmdString.WriteString(" ")
//  }
//  return cmdString.String()
//}

