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
	flags *Flags
	cwd string
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

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Cmd{
		cwd: cwd,
	}
}

func(cmd *Cmd) SetFlags(f *Flags) *Cmd {
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
		//log.Fatal("Command finished with error: %v\n", cmd.exec.String())
		fmt.Printf("%v\n", stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		return stdout.Bytes()
	}

	//fmt.Printf("%+V\n", string(cmd.Media.json))
	//fmt.Println(cmd.exec.String())
	//cmd.cmdArgs = []string{}
	return nil
}


func(cmd *Cmd) Show(action, input string) *Cmd {
	switch action {
	case "json":
		NewMedia(input).JsonMeta().Print()
		//fmt.Printf("%+V\n", string(cmd.Media.GetJsonMeta()))
	case "flags":
		//fmt.Printf("%+v\n", cmd.Flags)
	case "args":
		//fmt.Printf("%+v\n", cmd.args)
	case "meta":
		m := NewMedia(input).JsonMeta().Unmarshal()
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

func(c *Cmd) Extract(input string) {
	ffmpeg := NewFFmpegCmd(input).SetFlags(c.flags)

	switch {
	case c.flags.ChapSwitch:
		ffmpeg.media.FFmetaChapsToCue()
		return
	case c.flags.CoverSwitch:
		fmt.Println("cover")
		ffmpeg.AudioCodec = "an"
		ffmpeg.VideoCodec = "copy"
		ffmpeg.Output = "cover"
		ffmpeg.Ext = ".jpg"
	case c.flags.MetaSwitch:
		ffmpeg.PostInput = append(ffmpeg.PostInput, newMapArg("f", "ffmetadata"))
		ffmpeg.AudioCodec = "none"
		ffmpeg.VideoCodec = "none"
		ffmpeg.Output = "ffmeta"
		ffmpeg.Ext = ".ini"
	}
	ffmpeg.Parse()
	fmt.Println(ffmpeg.String())
	//c.Run()
}

