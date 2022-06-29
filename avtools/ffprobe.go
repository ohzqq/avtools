package avtools

import (
	"fmt"
	"os"
	"os/exec"
	//"bytes"
	//"log"
	//"strings"
	//"path/filepath"
)
var _ = fmt.Printf

type ffprobeCmd struct {
	Input string
	media *Media
	args cmdArgs
	exec *exec.Cmd
	tmpFile *os.File
	ffprobeArgs
}

type ffprobeArgs struct {
	LogLevel string
	pretty bool
	streams string
	entries string
	showChaps bool
	format string
}

func NewFFprobeCmd(i string) *ffprobeCmd {
	return &ffprobeCmd{
		Input: i,
		media: NewMedia(i),
	}
}

func(cmd *ffprobeCmd) String() string {
	return cmd.exec.String()
}

func(cmd *ffprobeCmd) Parse() Cmd {
	if log := Cfg().GetDefault("loglevel"); log != "" {
		cmd.args.Append("-v", log)
	}

	if cmd.pretty {
		cmd.args.Append("-pretty")
	}

	if stream := cmd.streams; stream != "" {
		cmd.args.Append("-select_streams", stream)
	}

	if entries := cmd.entries; entries != "" {
		cmd.args.Append("-show_entries", entries)
	}

	if cmd.showChaps {
		cmd.args.Append("-show_chapters")
	}

	cmd.args.Append("-of")
	switch f := cmd.format; f {
	default:
		fallthrough
	case "":
		fallthrough
	case "plain":
		cmd.args.Append("default=noprint_wrappers=1:nokey=1")
	case "json":
		cmd.args.Append("json=c=1")
	}

	cmd.args.Append(cmd.Input)

	return Cmd{exec: exec.Command("ffprobe", cmd.args.args...)}
}


