package avtools

import (
	//"path/filepath"
	//"log"
	"fmt"
	//"os"
	//"strings"
	//"os/exec"

	//"github.com/alessio/shellescape"
)
var _ = fmt.Printf

type Args struct {
	PreInput flagArgs
	PostInput flagArgs
	VideoCodec string
	VideoParams flagArgs
	VideoFilters []string
	AudioCodec string
	AudioParams flagArgs
	AudioFilters []string
	FilterComplex []string
	MiscParams []string
	Name string
	Profile string
	LogLevel string
	Output string
	Padding string
	Overwrite bool
	Extension string
	Start string
	End string
	Input string
	ChapNo int
	MetaFlag bool
	CoverFlag bool
	CueFlag bool
	ChapFlag bool
	Verbose bool
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

func NewArgs() *Args {
	return &Args{
		Output: Cfg().defaults.Output,
		LogLevel: Cfg().defaults.LogLevel,
		Overwrite: Cfg().defaults.Overwrite,
		Profile: Cfg().defaults.Profile,
		Padding: Cfg().defaults.Padding,
	}
}

type flagArgs []map[string]string
