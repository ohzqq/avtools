package avtools

import (
	//"path/filepath"
	//"log"
	"fmt"
	//"os"
	"strings"
	//"strconv"
	//"os/exec"

	//"github.com/alessio/shellescape"
)
var _ = fmt.Printf

type Args struct {
	Flags
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

func NewArgs() *Args {
	return &Args{
		Flags: Flags{Profile: "default"},
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
