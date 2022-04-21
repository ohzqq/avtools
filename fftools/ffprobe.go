package fftools

import (
	"log"
	"os"
	"fmt"
	"os/exec"
	"bytes"
)

type MediaMeta struct {
	Streams Streams
	Format Format
}

type Streams []Stream

type Stream struct {
	CodecName string
	CodecType string
}

type Format struct {
	Filename string
	StartTime string
	Duration string
	Size string
	BitRate string
	Tags Tags
}

type Tags struct {
	Title string
	Artist string
	Composer string
	Album string
	Comment string
	Genre string
}

type FFProbeCmd struct {
	cmd *exec.Cmd
	Input string
	tmpFile *os.File
	args ProbeArgs
}

func NewFFProbeCmd() *FFProbeCmd {
	ff := FFProbeCmd{}
	ff.cmd = exec.Command("ffprobe", "-hide_banner")
	return &ff
}

func AllJsonMeta() *FFProbeCmd {
	cmd := NewFFProbeCmd()
	cmd.In("/mnt/roar/audiobook/tmp/scrape/Palm_Island_1_-_Perfect_Ten_-_K.M._Neuhold/Palm_Island_01_-_Perfect_Ten.m4b")
	cmd.Args().
		Entries("format=filename,start_time,duration,size,bit_rate,format_tags:stream=codec_type,codec_name:format_tags").
		Chapters(true).
		Verbosity("error").
		Format("json")
	fmt.Printf("%v\n", cmd.Cmd().String())
	return cmd
}

func (ff *FFProbeCmd) Run() {
	//defer os.Remove(ff.tmpFile.Name())

	cmd := ff.Cmd()

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v\n", err)
		fmt.Printf("%v\n", stderr.String())
	}
	fmt.Printf("%v\n", stdout.String())
}

func (ff *FFProbeCmd) Cmd() *exec.Cmd {
	argOrder := []string{"Verbosity", "Streams", "Entries", "Chapters", "Params", "Format", "Input"}

	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.Verbosity()
		//case "Input":
		//  ff.pushInput()
		case "Streams":
			ff.Streams()
		case "Entries":
			ff.Entries()
		case "Chapters":
			ff.Chapters()
		//case "Params":
		//  ff.Params()
		case "Format":
			ff.Format()
		}
	}
	return ff.cmd
}

func (ff *FFProbeCmd) push(arg string) {
	ff.cmd.Args = append(ff.cmd.Args, arg)
}

func (ff *FFProbeCmd) Verbosity() {
	ff.push("-v")
	if ff.args.verbosity != "" {
		ff.push(ff.args.verbosity)
	} else {
		ff.push("fatal")
	}
}

func (ff *FFProbeCmd) In(input string) {
	ff.push(input)
	//ff.Input = input
}

func (ff *FFProbeCmd) pushInput() {
	ff.push(ff.Input)
}

func (ff *FFProbeCmd) Pretty() {
	if ff.args.pretty {
		ff.push("-pretty")
	}
}

func (ff *FFProbeCmd) Chapters() {
	if ff.args.chapters {
		ff.push("-show_chapters")
	}
}

//func (ff *FFProbeCmd) Params() {
//  if ff.args.params != "" {
//    ff.push("-c:v")
//    ff.push(ff.args.params)
//  }
//}

func (ff *FFProbeCmd) Entries() {
	if ff.args.entries != "" {
		ff.push("-show_entries")
		ff.push(ff.args.entries)
	}
}

func (ff *FFProbeCmd) Streams() {
	if ff.args.streams != "" {
		ff.push("-select_streams")
		ff.push(ff.args.streams)
	}
}

func (ff *FFProbeCmd) Format() {
	ff.push("-of")
	switch f := ff.args.format; f {
	default:
		fallthrough
	case "":
		fallthrough
	case "plain":
		ff.push("default=noprint_wrappers=1:nokey=1")
	case "json":
		ff.push("json=c=1")
	}
}

func (ff *FFProbeCmd) Args() *ProbeArgs {
	ff.args = ProbeArgs{}
	return &ff.args
}

type ProbeArgs struct {
	pretty bool
	streams string
	entries string
	chapters bool
	format string
	//params string
	verbosity string
}

func (a *ProbeArgs) Pretty(arg bool) *ProbeArgs {
	a.pretty = arg
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

func (a *ProbeArgs) Chapters(arg bool) *ProbeArgs {
	a.chapters = arg
	return a
}

func (a *ProbeArgs) Format(arg string) *ProbeArgs {
	a.format = arg
	return a
}

//func (a *ProbeArgs) Params(arg string) *ProbeArgs {
//  a.params = arg
//  return a
//}

func (a *ProbeArgs) Verbosity(arg string) *ProbeArgs {
	a.verbosity = arg
	return a
}

