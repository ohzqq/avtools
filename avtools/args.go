package avtools

import (
	"path/filepath"
	"log"
	"fmt"
	//"os"
	//"strings"
	//"os/exec"

	//"github.com/alessio/shellescape"
)
var _ = fmt.Printf

type CmdArgs struct {
	PreInput flagArgs
	PostInput flagArgs
	VideoCodec string
	VideoParams flagArgs
	VideoFilters string
	AudioCodec string
	AudioParams flagArgs
	AudioFilters string
	FilterComplex []string
	MiscParams []string
	Name string
	Profile string
	LogLevel string
	Output string
	Padding string
	Extension string
	Start string
	ChapNo int
	Cover bool
	MetaFlag bool
	CueFlag bool
	ChapFlag bool
	Verbose bool
	Overwrite bool
	CoverFile string
	MetaFile string
	CueFile string
	num int
	pretty bool
	streams string
	entries string
	chapters bool
	format string
}

func NewArgs() *CmdArgs {
	return &CmdArgs{
		Output: Cfg().defaults.Output,
		LogLevel: Cfg().defaults.LogLevel,
		Overwrite: Cfg().defaults.Overwrite,
		Profile: Cfg().defaults.Profile,
		Padding: Cfg().defaults.Padding,
	}
}

type flagArgs []map[string]string

func (f flagArgs) Split() []string {
	var args []string
	for _, flArg := range f {
		for flag, arg := range flArg {
			flag = "-" + flag
			args = append(args, flag, arg)
		}
	}
	return args
}

func (a *CmdArgs) Cover(s string) *CmdArgs {
	path, err := filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}
	a.AlbumArt = path
	return a
}

func (a *CmdArgs) Cue(s string) *CmdArgs {
	a.CueSheet = s
	return a
}

func (a *CmdArgs) Num(n int) *CmdArgs {
	a.num = n
	return a
}

func (a *CmdArgs) Meta(s string) *CmdArgs {
	path, err := filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}
	a.Metadata = path
	return a
}

func (a *CmdArgs) LogLevel(s string) *CmdArgs {
	a.Verbosity = s
	return a
}

func (a *CmdArgs) Pre(k, v string) *CmdArgs {
	a.PreInput = append(a.PreInput, map[string]string{k: v})
	return a
}

func (a *CmdArgs) OverWrite() *CmdArgs {
	a.Overwrite = true
	return a
}

func (a *CmdArgs) Verbose() *CmdArgs {
	a.verbose = true
	return a
}

func (a *CmdArgs) Post(k, v string) *CmdArgs {
	a.PostInput = append(a.PostInput, map[string]string{k: v})
	return a
}

func (a *CmdArgs) VCodec(s string) *CmdArgs {
	a.VideoCodec = s
	return a
}

func (a *CmdArgs) VParams(k, v string) *CmdArgs {
	a.VideoParams = append(a.VideoParams, map[string]string{k: v})
	return a
}

func (a *CmdArgs) VFilters(s string) *CmdArgs {
	a.VideoFilters = s
	return a
}

func (a *CmdArgs) Filters(f []string) *CmdArgs {
	a.FilterComplex = f
	return a
}

func (a *CmdArgs) Params(p []string) *CmdArgs {
	a.MiscParams = p
	return a
}

func (a *CmdArgs) ACodec(s string) *CmdArgs {
	a.AudioCodec = s
	return a
}

func (a *CmdArgs) AParams(k, v string) *CmdArgs {
	a.AudioParams = append(a.AudioParams, map[string]string{k: v})
	return a
}

func (a *CmdArgs) AFilters(s string) *CmdArgs {
	a.AudioFilters = s
	return a
}

func (a *CmdArgs) Pad(s string) *CmdArgs {
	a.Padding = s
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

func (a *CmdArgs) Pretty() *CmdArgs {
	a.pretty = true
	return a
}

func (a *CmdArgs) Streams(arg string) *CmdArgs {
	a.streams = arg
	return a
}

func (a *CmdArgs) Entries(arg string) *CmdArgs {
	a.entries = arg
	return a
}

func (a *CmdArgs) Chapters() *CmdArgs {
	a.chapters = true
	return a
}

func (a *CmdArgs) Format(arg string) *CmdArgs {
	a.format = arg
	return a
}

func Mp3CoverArgs() []string {
	return []string{
		"-id3v2_version",
		"3",
		"-metadata:s:v",
		"title='Album cover'",
		"-metadata:s:v",
		"comment='Cover (front)'",
	}
}
