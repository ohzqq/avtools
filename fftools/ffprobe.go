package fftools

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

type FFProbeCmd struct {
	cmd *exec.Cmd
	Input string
	args ProbeArgs
}

type ProbeArgs struct {
	Pretty bool
	Streams string
	Entries string
	Chapters bool
	Format string
	Params string
	Verbosity string
}

func NewFFProbeCmd() *FFProbeCmd {
	ff := FFmpegCmd{}
	ff.cmd = exec.Command("ffprobe", "-hide_banner")
	return &ff
}

func (ff *FFProbeCmd) Args() *ProbeArgs {
	ff.args = ProbeArgs{}
	ff.args.VCodec("copy")
	ff.args.ACodec("copy")
	return &ff.args
}

