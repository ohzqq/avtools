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

func (f Flag) Parse() *Cmd {
	c := NewCmd()

	if f.Args.Profile != "" {
		c.Profile = Cfg().GetProfile(f.Args.Profile)
	}

	if f.Args.Start != "" {
		c.Start = f.Args.Start
	}

	if f.Args.End != "" {
		c.End = f.Args.End
	}

	if f.Args.Output != "" {
		c.Output = file.New(f.Args.Output)
	}

	if f.Args.Input != "" {
		c.Input = file.New(f.Args.Input)
		c.Media = media.NewMedia(f.Args.Input)
	}

	if f.Args.Cover != "" {
		c.Cover = file.New(f.Args.Cover)
		c.Media.AddFile("cover", f.Args.Cover)
	}

	if f.Args.Meta != "" {
		c.Meta = file.New(f.Args.Meta)
		c.Media.SetFFmeta(f.Args.Meta)
	}

	if f.Args.Cue != "" {
		c.Cue = file.New(f.Args.Cue)
		c.Media.SetCue(f.Args.Cue)
	}

	if f.Args.Input != "" {
		c.Media.SetMeta()
	}

	c.PadOutput = Cfg().Defaults.HasPadding()
	c.Padding = Cfg().Defaults.Padding
	c.Num = 1

	if f.Args.ChapNo != 0 {
		c.ChapNo = f.Args.ChapNo
	}

	return c
}
