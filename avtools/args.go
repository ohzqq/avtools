package avtools

import (
	"fmt"
	"strings"
)

type Args struct {
	Options
	Input string
	PreInput mapArgs
	PostInput mapArgs
	VideoCodec string
	VideoParams mapArgs
	VideoFilters stringArgs
	AudioCodec string
	AudioParams mapArgs
	AudioFilters stringArgs
	FilterComplex stringArgs
	MiscParams stringArgs
	LogLevel string
	Name string
	Padding string
	Ext string
	num int
	pretty bool
	streams string
	entries string
	showChaps bool
	format string
}

type Options struct {
	Overwrite bool
	Profile string
	Start string
	End string
	Output string
	ChapNo int
	MetaSwitch bool
	CoverSwitch bool
	CueSwitch bool
	ChapSwitch bool
	Verbose bool
	CoverFile string
	MetaFile string
	CueFile string
}

func NewArgs() *Args {
	return &Args{
		Options: Options{Profile: "default"},
	}
}

type cmdArgs struct {
	args []string
}

func(arg *cmdArgs) Append(args ...string) {
	arg.args = append(arg.args, args...)
}

type mapArgs []map[string]string

func newMapArg(k, v string) map[string]string {
	return map[string]string{k: v}
}

func(a *Args) AppendMapArg(key, flag, value string) {
	mapArg := newMapArg(flag, value)
	switch key {
	case "pre":
		a.PreInput = append(a.PreInput, mapArg)
	case "post":
		a.PostInput = append(a.PostInput, mapArg)
	case "videoParams":
		a.VideoParams = append(a.VideoParams, mapArg)
	case "audioParams":
		a.AudioParams = append(a.AudioParams, mapArg)
	}
}

func(m mapArgs) Split() []string {
	var args []string
	for _, flArg := range m {
		for flag, arg := range flArg {
			args = append(args, "-" + flag, arg)
		}
	}
	return args
}

type stringArgs []string

func(s stringArgs) Join() string {
	return strings.Join(s, ",")
}
