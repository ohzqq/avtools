package tool

import (
	"github.com/ohzqq/avtools/file"
	"github.com/ohzqq/avtools/media"
)

type Flag struct {
	Args ArgFlag
	Bool BoolFlag
}

type BoolFlag struct {
	Meta      bool
	Cover     bool
	Cue       bool
	Chap      bool
	Json      bool
	Overwrite bool
	Verbose   bool
}

type ArgFlag struct {
	Profile string
	Start   string
	End     string
	Output  string
	ChapNo  int
	Input   string
	Cover   string
	Meta    string
	Cue     string
}

type Args struct {
	Profile   Profile
	Start     string
	End       string
	output    file.File
	Input     file.File
	Cover     file.File
	Meta      file.File
	Cue       file.File
	Media     *media.Media
	Num       int
	PadOutput bool
	Padding   string
	ChapNo    int
}

func (a Args) Output() string {
	if a.PadOutput {
		return a.output.Pad(a.Num)
	}
	return a.output.Name
}

//func (f Flag) Parse() Args {
//  return Args{
//    Profile:   f.Args.GetProfile(),
//    Start:     f.Args.Start,
//    End:       f.Args.End,
//    output:    file.New(f.Args.Output),
//    Input:     file.New(f.Args.Input),
//    Cover:     file.New(f.Args.Cover),
//    Meta:      file.New(f.Args.Meta),
//    Cue:       file.New(f.Args.Cue),
//    Media:     f.Media(),
//    PadOutput: Cfg().Defaults.HasPadding(),
//    Padding:   Cfg().Defaults.Padding,
//    Num:       1,
//    ChapNo:    f.Args.ChapNo,
//  }
//}

//func (f Flag) Media() *media.Media {
//  var m *media.Media
//  if f.Args.HasInput() {
//    m = media.NewMedia(f.Args.Input)

//    if f.Args.HasMeta() {
//      m.SetFFmeta(f.Args.Meta)
//    }

//    if f.Args.HasCue() {
//      m.SetCue(f.Args.Cue)
//    }

//    if f.Args.HasCover() {
//      m.AddFile("cover", f.Args.Cover)
//    }

//    m.SetMeta()
//  }

//  return m
//}

func (f Cmd) HasCover() bool {
	return f.Cover.Abs != ""
}

func (f Cmd) HasChapNo() bool {
	return f.ChapNo != 0
}

func (f Cmd) HasCue() bool {
	return f.Cue.Abs != ""
}

func (f Cmd) HasMeta() bool {
	return f.Meta.Abs != ""
}

func (f Cmd) HasStart() bool {
	return f.Start != ""
}

func (f Cmd) HasEnd() bool {
	return f.End != ""
}

func (f Cmd) HasInput() bool {
	return f.Input.Abs != ""
}

func (f Cmd) HasOutput() bool {
	return f.Output.Abs != ""
}
