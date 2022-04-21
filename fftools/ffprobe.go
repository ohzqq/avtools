package fftools

import (
	"log"
	"os"
	"fmt"
	"os/exec"
	"bytes"
	"encoding/json"
)

type MediaMeta struct {
	Chapters Chapters
	Streams Streams
	Format Format
}

type Streams []Stream

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	StartTime string `json:"start_time"`
	Duration string
	Size string
	BitRate string `json:"bit_rate"`
	Tags Tags
}

type Tags struct {
	Title string `json:"title"`
	Artist string `json:"artist"`
	Composer string `json:"composer"`
	Album string `json:"album"`
	Comment string `json:"comment"`
	Genre string `json:"genre"`
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

func (ff *FFProbeCmd) AllJsonMeta(input string) *MediaMeta {
	ff.In(input)
	ff.Args().
		Entries("format=filename,start_time,duration,size,bit_rate,format_tags:stream=codec_type,codec_name:format_tags").
		Chapters().
		Verbosity("error").
		Format("json")

	m := ff.Run()
	var meta MediaMeta
	err := json.Unmarshal(m, &meta)
	if err != nil { fmt.Println("help")}
	return &meta
}

func (ff *FFProbeCmd) Run() []byte {
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
	return stdout.Bytes()
}

func (ff *FFProbeCmd) Cmd() *exec.Cmd {
	argOrder := []string{"Verbosity", "Streams", "Entries", "Chapters", "Pretty", "Format", "Input"}

	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.Verbosity()
		case "Pretty":
			ff.Pretty()
		case "Streams":
			ff.Streams()
		case "Entries":
			ff.Entries()
		case "Chapters":
			ff.Chapters()
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

