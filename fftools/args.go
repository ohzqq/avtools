package fftools

import (
	//"os"
	//"strings"
	//"os/exec"

	//"github.com/alessio/shellescape"
)

type CmdArgs struct {
	Pre string
	Post string
	VideoCodec string
	VideoParams string
	VideoFilters string
	AudioCodec string
	AudioParams string
	AudioFilters string
	FilterComplex string
	Verbosity string
	Output string
	Padding bool
	Ext string
	Overwrite bool
	CueSheet string
	AlbumCover string
	Metadata string
}

func (a *CmdArgs) Cover(s string) *CmdArgs {
	a.AlbumCover = s
	return a
}

func (a *CmdArgs) Cue(s string) *CmdArgs {
	a.CueSheet = s
	return a
}

func (a *CmdArgs) Meta(s string) *CmdArgs {
	a.Metadata = s
	return a
}

func (a *CmdArgs) LogLevel(s string) *CmdArgs {
	a.Verbosity = s
	return a
}

func (a *CmdArgs) PreInput(s string) *CmdArgs {
	a.Pre = s
	return a
}

func (a *CmdArgs) OverWrite(over bool) *CmdArgs {
	a.Overwrite = over
	return a
}

func (a *CmdArgs) PostInput(s string) *CmdArgs {
	a.Post = s
	return a
}

func (a *CmdArgs) VCodec(s string) *CmdArgs {
	a.VideoCodec = s
	return a
}

func (a *CmdArgs) VParams(s string) *CmdArgs {
	a.VideoParams = s
	return a
}

func (a *CmdArgs) VFilters(s string) *CmdArgs {
	a.VideoFilters = s
	return a
}

func (a *CmdArgs) Filter(s string) *CmdArgs {
	a.FilterComplex = s
	return a
}

func (a *CmdArgs) ACodec(s string) *CmdArgs {
	a.AudioCodec = s
	return a
}

func (a *CmdArgs) AParams(s string) *CmdArgs {
	a.AudioParams = s
	return a
}

func (a *CmdArgs) AFilters(s string) *CmdArgs {
	a.AudioFilters = s
	return a
}

func (a *CmdArgs) Out(s string)  {
	a.Output = s
	//return a
}
