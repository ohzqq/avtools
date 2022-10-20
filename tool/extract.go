package tool

import (
	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/file"
)

type ExtractCmd struct {
	*Cmd
}

func Extract() *ExtractCmd {
	return &ExtractCmd{Cmd: NewCmd()}
}

func (e *ExtractCmd) Parse() *Cmd {
	if e.flag.Bool.Cover && e.Media().HasEmbeddedCover() {
		out := "cover" + e.Media().EmbeddedCoverExt()
		ff := e.FFmpeg()
		ff.AN().CV("copy").Output(out)
		e.Add(ff)
	}

	if e.flag.Bool.Meta {
		//ff := ffmpeg.New()
		//ff.AppendPostInput("f", "ffmetadata").Output("ffmeta.ini").Input(e.Args.Input.Abs).Overwrite()
		ff := ExtractFFmeta(e.Args.Input.Abs)
		e.Add(ff)
	}

	if e.flag.Bool.Cue {
		e.Media().Meta.Chapters.File = "chapters.cue"
		e.Media().Meta.Chapters.Write()
	}

	return e.Cmd
}

func ExtractFFmeta(input string) *ffmpeg.Cmd {
	i := file.New(input)
	ff := ffmpeg.New()
	ff.AppendPostInput("f", "ffmetadata").
		Output("ffmeta.ini").
		Input(i.Abs).
		Overwrite()
	return ff
}
