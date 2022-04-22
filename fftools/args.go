package fftools

import (
	//"os"
	//"strings"
	//"os/exec"

	//"github.com/alessio/shellescape"
)

type CmdArgs struct {
	Pre flagArgs
	Post flagArgs
	VideoCodec string
	VideoParams flagArgs
	VideoFilters string
	AudioCodec string
	AudioParams flagArgs
	AudioFilters string
	FilterComplex string
	Verbosity string
	Output string
	Padding bool
	Extension string
	Overwrite bool
	CueSheet string
	AlbumCover string
	Metadata string
}

type flagArgs map[string]string

func newFlagArg(flag, arg string) flagArgs {
	return flagArgs{flag: arg}
}

func (f flagArgs) Split() []string {
	var args []string
	for flag, arg := range f {
		flag = "-" + flag
		args = append(args, flag, arg)
	}
	return args
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

func (a *CmdArgs) PreInput(s flagArgs) *CmdArgs {
	a.Pre = s
	return a
}

func (a *CmdArgs) OverWrite(over bool) *CmdArgs {
	a.Overwrite = over
	return a
}

func (a *CmdArgs) PostInput(s flagArgs) *CmdArgs {
	a.Post = s
	return a
}

func (a *CmdArgs) VCodec(s string) *CmdArgs {
	a.VideoCodec = s
	return a
}

func (a *CmdArgs) VParams(f flagArgs) *CmdArgs {
	a.VideoParams = f
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

func (a *CmdArgs) AParams(s flagArgs) *CmdArgs {
	a.AudioParams = s
	return a
}

func (a *CmdArgs) AFilters(s string) *CmdArgs {
	a.AudioFilters = s
	return a
}

func (a *CmdArgs) Ext(s string) *CmdArgs {
	a.Extension = s
	return a
}

func (a *CmdArgs) Out(s string) *CmdArgs {
	a.Output = s
	return a
}
