package avtools

import (
	"log"
	"os"
	"fmt"
	"os/exec"
	"bytes"
)

type FFprobeCmd struct {
	cmd *exec.Cmd
	Input string
	tmpFile *os.File
	verbose bool
	args ProbeArgs
}

func NewFFprobeCmd() *FFprobeCmd {
	return &FFprobeCmd{cmd: exec.Command("ffprobe", "-hide_banner")}
}

func (ff *FFprobeCmd) In(input string) {
	ff.Input = input
}

func (ff *FFprobeCmd) Args() *ProbeArgs {
	ff.args = ProbeArgs{}
	return &ff.args
}

func (ff *FFprobeCmd) Run() []byte {
	cmd := ff.buildCmd()

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v\n", cmd.String())
		fmt.Printf("%v\n", stderr.String())
	}

	if ff.verbose {
		fmt.Println(cmd.String())
		fmt.Println(stdout.String())
	}

	return stdout.Bytes()
}

func(ff *FFprobeCmd) String() string {
	return ff.buildCmd().String()
}

func (ff *FFprobeCmd) buildCmd() *exec.Cmd {
	argOrder := []string{"Verbosity", "Streams", "Entries", "Chapters", "Pretty", "Format", "Input"}

	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.cmd.Args = append(ff.cmd.Args, "-v")
			if ff.args.verbosity != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.args.verbosity)
			} else {
				ff.cmd.Args = append(ff.cmd.Args, "fatal")
			}
		case "Pretty":
			if ff.args.pretty {
				ff.cmd.Args = append(ff.cmd.Args, "-pretty")
			}
		case "Streams":
			if ff.args.streams != "" {
				ff.cmd.Args = append(ff.cmd.Args, "-select_streams", ff.args.streams)
			}
		case "Entries":
			if ff.args.entries != "" {
				ff.cmd.Args = append(ff.cmd.Args, "-show_entries", ff.args.entries)
			}
		case "Chapters":
			if ff.args.chapters {
				ff.cmd.Args = append(ff.cmd.Args, "-show_chapters")
			}
		case "Format":
			ff.cmd.Args = append(ff.cmd.Args, "-of")
			switch f := ff.args.format; f {
			default:
				fallthrough
			case "":
				fallthrough
			case "plain":
				ff.cmd.Args = append(ff.cmd.Args, "default=noprint_wrappers=1:nokey=1")
			case "json":
				ff.cmd.Args = append(ff.cmd.Args, "json=c=1")
			}
		case "Input":
			ff.cmd.Args = append(ff.cmd.Args, ff.Input)
		}
	}
	return ff.cmd
}

type ProbeArgs struct {
	pretty bool
	streams string
	entries string
	chapters bool
	format string
	verbosity string
}

func (a *ProbeArgs) Pretty() *ProbeArgs {
	a.pretty = true
	return a
}

func (a *ProbeArgs) Streams(arg string) *ProbeArgs {
	a.streams = arg
	return a
}

func (a *ProbeArgs) Entries(arg string) *ProbeArgs {
	a.entries = arg
	return a
}

func (a *ProbeArgs) Chapters() *ProbeArgs {
	a.chapters = true
	return a
}

func (a *ProbeArgs) Format(arg string) *ProbeArgs {
	a.format = arg
	return a
}

func (a *ProbeArgs) Verbosity(arg string) *ProbeArgs {
	a.verbosity = arg
	return a
}

