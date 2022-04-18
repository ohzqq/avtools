package fftools

import (
	//"os"
	//"strings"
	"os/exec"
)

type CmdArgs map[string]string

var ArgOrder = []string{"Verbosity", "Pre", "Input", "Post", "VideoCodec", "VideoParams", "VideoFilters", "AudioCodec", "AudioParams", "AudioFilters", "FilterComplex", "Output"}

type FFmpegCmd struct {
	cmd *exec.Cmd
	Input string
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
	ff.Arguments.SetVerbosity()
	return ff
}

func (c CmdArgs) SetVerbosity() {
	if Cfg.Defaults.IsSet("Verbosity") {
		c["Verbosity"] = "-loglevel " + Cfg.Defaults.GetString("Verbosity")
	}
}

func (c CmdArgs) SetInput(arg string) {
	c["Input"] = arg
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	for _, arg := range ArgOrder {
		ff.cmd.Args = append(ff.cmd.Args, ff.Arguments[arg])
	}
	return ff.cmd
}

func (ff *FFmpegCmd) PreInput(arg string) {
	ff.Arguments["Pre"] = arg
}

func (ff *FFmpegCmd) PostInput(arg string) {
	ff.Arguments["Post"] = arg
}

func (ff *FFmpegCmd) VC(arg string) {
	ff.Arguments["VideoCodec"] = arg
}

func (ff *FFmpegCmd) VP(arg string) {
	ff.Arguments["VideoParams"] = arg
}

func (ff *FFmpegCmd) VF(arg string) {
	ff.Arguments["VideoFilters"] = arg
}

func (ff *FFmpegCmd) FilterComplex(arg string) {
	ff.Arguments["FilterComplex"] = arg
}

func (ff *FFmpegCmd) AC(arg string) {
	ff.Arguments["AudioCodec"] = arg
}

func (ff *FFmpegCmd) AP(arg string) {
	ff.Arguments["AudioParams"] = arg
}

func (ff *FFmpegCmd) AF(arg string) {
	ff.Arguments["AudioFilters"] = arg
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

