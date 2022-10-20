package tool

import "github.com/ohzqq/avtools/ffmpeg"

type ExtractCmd struct {
	*Cmd
}

func Extract() *ExtractCmd {
	return &ExtractCmd{Cmd: NewCmd()}
}

func (e *ExtractCmd) Parse() *Cmd {
	if e.flag.Bool.Cover && e.Args.Media.HasEmbeddedCover() {
		out := "cover" + e.Args.Media.EmbeddedCoverExt()
		ff := ffmpeg.New()
		ff.AN().CV("copy").Output(out)
		e.Add(ff)
	}

	if e.flag.Bool.Meta {
		ff := ffmpeg.New()
		filter := ffmpeg.NewFilter("ffmetadata")
		ff.Filter(filter).Output("ffmeta.ini")
		e.Add(ff)
	}

	if e.flag.Bool.Cue {
		e.Media.Meta.Chapters.File = "chapters.cue"
		e.Media.Meta.Chapters.Write()
	}

	return e.Cmd
}
