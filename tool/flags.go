package media

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

func (f Flag) Media() *Media {
	var media *Media
	if f.Args.HasInput() {
		media = NewMedia(f.Args.Input)

		if f.Args.HasMeta() {
			media.SetFFmeta(f.Args.Meta)
		}

		if f.Args.HasCue() {
			media.SetCue(f.Args.Cue)
		}

		if f.Args.HasCover() {
			media.AddFile("cover", f.Args.Cover)
		}

		media.SetMeta()
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
