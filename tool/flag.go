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
}

func (a Args) Output() string {
	if a.PadOutput {
		return a.output.Pad(a.Num)
	}
	return a.output.Name
}

func (f Flag) Parse() Args {
	return Args{
		Profile:   f.Args.GetProfile(),
		Start:     f.Args.Start,
		End:       f.Args.End,
		output:    file.New(f.Args.Output),
		Input:     file.New(f.Args.Input),
		Cover:     file.New(f.Args.Cover),
		Meta:      file.New(f.Args.Meta),
		Cue:       file.New(f.Args.Cue),
		Media:     f.Media(),
		PadOutput: Cfg().Defaults.HasPadding(),
		Padding:   Cfg().Defaults.Padding,
		Num:       1,
	}
}

func (f Flag) Media() *media.Media {
	var m *media.Media
	if f.Args.HasInput() {
		m = media.NewMedia(f.Args.Input)

		if f.Args.HasMeta() {
			m.SetFFmeta(f.Args.Meta)
		}

		if f.Args.HasCue() {
			m.SetCue(f.Args.Cue)
		}

		if f.Args.HasCover() {
			m.AddFile("cover", f.Args.Cover)
		}

		m.SetMeta()
	}

	return m
}

func (f ArgFlag) HasCover() bool {
	return f.Cover != ""
}

func (f ArgFlag) HasCue() bool {
	return f.Cue != ""
}

func (f ArgFlag) HasMeta() bool {
	return f.Meta != ""
}

func (f ArgFlag) HasProfile() bool {
	return f.Profile != ""
}

func (f ArgFlag) GetProfile() Profile {
	if f.HasProfile() {
		return Cfg().GetProfile(f.Profile)
	}
	return Cfg().GetProfile("default")
}

func (f ArgFlag) HasStart() bool {
	return f.Start != ""
}

func (f ArgFlag) HasEnd() bool {
	return f.End != ""
}

func (f ArgFlag) HasInput() bool {
	return f.Input != ""
}

func (f ArgFlag) HasOutput() bool {
	return f.Output != ""
}
