package avtools

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"log"
	//"strings"
	//"path/filepath"
)
var _ = fmt.Printf

type Cmd struct {
	input string
	flags *Flags
	cwd string
	ffmpeg *ffmpegCmd
	exec *exec.Cmd
	tmpFile *os.File
}

type Flags struct {
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

func NewCmd(i string) *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Cmd{
		cwd: cwd,
		input: i,
	}
}

func(cmd *Cmd) FFmpeg() *ffmpegCmd {
	cmd.ffmpeg = &ffmpegCmd{
		flags: cmd.flags,
		Args: Cfg().GetProfile(cmd.flags.Profile),
		media: NewMedia(cmd.input),
	}
	return cmd.ffmpeg.ParseFlags()
}

func(cmd *Cmd) Options(f *Flags) *Cmd {
	cmd.flags = f
	return cmd
}

func(cmd Cmd) Run() []byte {
	if cmd.tmpFile != nil {
		defer os.Remove(cmd.tmpFile.Name())
	}

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.exec.Stderr = &stderr
	cmd.exec.Stdout = &stdout

	err := cmd.exec.Run()
	if err != nil {
		log.Fatal("Command finished with error: %v\n", cmd.exec.String())
		fmt.Printf("%v\n", stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		return stdout.Bytes()
	}

	if cmd.flags.Verbose {
		fmt.Println(cmd.String())
	}
	return nil
}

func(cmd Cmd) String() string {
	return cmd.exec.String()
}

func(cmd *Cmd) Show(action string) *Cmd {
	media := NewMedia(cmd.input)
	switch action {
	case "json":
		media.JsonMeta().Print()
		//fmt.Printf("%+V\n", string(cmd.Media.GetJsonMeta()))
	case "flags":
		//fmt.Printf("%+v\n", cmd.Flags)
	case "args":
		//fmt.Printf("%+v\n", cmd.args)
	case "meta":
		m := media.JsonMeta().Unmarshal()
		fmt.Printf("%+V\n", m.Meta)
	case "cmd":
		//m := NewMedia(input).JsonMeta().Unmarshal()
		//fmt.Printf("%+V\n", m.Meta)
		//cmd.ffmpeg = true
		//cmd.ffprobe = true
		//fmt.Printf("%+v\n", Cfg().GetProfile(cmd.Flags.Profile))
		//fmt.Printf("%+v\n", cmd.exec.String())
	default:
		fmt.Printf("%+v\n", cmd)
	}
	return cmd
}

