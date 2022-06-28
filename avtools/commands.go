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
	Media *Media
	args *Args
	Flags *Flags
	Action string
	cmdArgs []string
	//FFmpegCmd *FFmpegCmd
	//FFprobeCmd *FFprobeCmd
	Input string
	Ext string
	cwd string
	exec *exec.Cmd
	tmpFile *os.File
	ffmpeg bool
	ffprobe bool
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

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Cmd{
		cwd: cwd,
		cmdArgs: []string{"-hide_banner"},
		Flags: &Flags{Profile: "default"},
	}
}

func(cmd *Cmd) Run() []byte {
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
		log.Printf("Command finished with error: %v\n", cmd.exec.String())
		fmt.Printf("%v\n", stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		return stdout.Bytes()
		//fmt.Printf("%v\n", stdout.String())
	}

	return nil
}

func(cmd *Cmd) ParseFlags() {
	cmd.args = Cfg().GetProfile(cmd.Flags.Profile)
	cmd.Media = NewMedia(cmd.Input)
	cmd.ParseJsonMeta()

	if meta := cmd.Flags.MetaFile; meta != "" {
		cmd.Media.SetMeta(LoadFFmetadataIni(meta))
	}

	if cue := cmd.Flags.CueFile; cue != "" {
		cmd.Media.SetChapters(LoadCueSheet(cue))
	}

	if y := cmd.Flags.Overwrite; y {
		cmd.args.Overwrite = y
	}

	if o := cmd.Flags.Output; o != "" {
		cmd.args.Name = o
	}

	if c := cmd.Flags.ChapNo; c  != 0 {
		cmd.args.num = c
	}

	//if e := cmd.Ext; e != "" {
	//  cmd.args.Ext = e
	//}
}

func(cmd *Cmd) Show() *Cmd {
	cmd.ParseFlags()
	switch cmd.Action {
	case "flags":
		fmt.Printf("%+v\n", cmd.Flags)
	case "args":
		fmt.Printf("%+v\n", cmd.args)
	case "meta":
		//cmd.Media.RenderFFChaps()
		fmt.Printf("%+V\n", cmd.Media.Meta)
	case "cmd":
		//cmd.ffmpeg = true
		//cmd.ffprobe = true
		//fmt.Printf("%+v\n", Cfg().GetProfile(cmd.Flags.Profile))
		fmt.Printf("%+v\n", cmd.exec.String())
	default:
		fmt.Printf("%+v\n", cmd)
	}
	return cmd
}

