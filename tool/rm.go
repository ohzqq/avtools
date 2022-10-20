package tool

type RmCmd struct {
	*Cmd
}

func Rm() *RmCmd {
	return &RmCmd{Cmd: NewCmd()}
}

func (rm *RmCmd) Parse() *Cmd {
	ff := rm.FFmpeg()
	if rm.flag.Bool.Cover && rm.Media().HasEmbeddedCover() {
		ff.VN()
	}

	if rm.flag.Bool.Meta {
		ff.AppendPostInput("map_metadata", "-1")
	}

	if rm.flag.Bool.Cue {
		ff.AppendPostInput("map_chapters", "-1")
	}

	rm.Add(ff)
	return rm.Cmd
}
