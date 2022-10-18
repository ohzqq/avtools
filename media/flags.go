package tool

import (
	"github.com/ohzqq/avtools/ffmpeg"
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

func (f Flag) FFmpegCmd() *ffmpeg.Cmd {
	cmd := ffmpeg.New()

	if f.Bool.Overwrite {
		cmd.AppendPreInput("y")
	}

	if f.Args.HasStart() {
		cmd.AppendPreInput("ss", f.Args.Start)
	}

	if f.Args.HasEnd() {
		cmd.AppendPreInput("to", f.Args.End)
	}

	if f.Args.HasInput() {
		cmd.Input(f.Args.Input)
	}

	if f.Args.HasMeta() {
		cmd.FFmeta(f.Args.Meta)
		//cmd.AppendPostInput("i", f.Args.Meta)
		//cmd.AppendPostInput("map_metadata", "1")
	}

	if f.Args.HasOutput() {
		cmd.Output(f.Args.Output)
	} else {
		cmd.Output(OutputFromInput(f.Args.Input).String())
	}

	return cmd
}

func (f Flag) Media() *Media {
	var media *Media
	if f.Args.HasInput() {
		media = NewMedia(f.Args.Input)
		if f.Args.HasMeta() {
			media.SetFile("ffmeta", f.Args.Meta)
		}
		if f.Args.HasCue() {
			media.SetFile("cue", f.Args.Cue)
		}
		if f.Args.HasCover() {
			media.SetFile("cover", f.Args.Cover)
		}
	}
	return media
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
