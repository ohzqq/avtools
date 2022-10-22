package tool

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/file"
	"github.com/ohzqq/avtools/media"
)

type Cmd struct {
	Profile   Profile
	Start     string
	End       string
	Output    file.File
	Input     file.File
	Cover     file.File
	Meta      file.File
	Cue       file.File
	Media     *media.Media
	Num       int
	PadOutput bool
	Padding   string
	ChapNo    int
	flag      BoolFlag
	flags     Flag
	isVerbose bool
	cwd       string
	Batch     []Command
	bin       string
	args      []string
	tmp       *os.File
}

type Command interface {
	Build() (*exec.Cmd, error)
	String() string
	ParseArgs() ([]string, error)
	Run() ([]byte, error)
}

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Cmd{
		Profile: Cfg().GetProfile("default"),
		cwd:     cwd,
	}
}

func (c *Cmd) SetInput(i string) *Cmd {
	c.Input = file.New(i)
	c.Media = media.NewMedia(i)
	return c
}

func (c Cmd) String() string {
	return strings.Join(c.args, " ")
}

//func (c Cmd) Media() *media.Media {
//  return c.Args.Media
//}

func (c Cmd) ParseArgs() ([]string, error) {
	return c.args, nil
}

func (c Cmd) Build() (*exec.Cmd, error) {
	cmd := exec.Command(c.bin, c.args...)
	return cmd, nil
}

func (c *Cmd) Bin(bin string) *Cmd {
	c.bin = bin
	return c
}

func (c *Cmd) SetArgs(args ...string) *Cmd {
	c.args = args
	return c
}

func (c *Cmd) Command(bin string, args []string) *Cmd {
	c.args = args
	c.bin = bin
	return c
}

func (c *Cmd) Add(cmd Command) *Cmd {
	c.Batch = append(c.Batch, cmd)
	return c
}

func (c *Cmd) Verbose() *Cmd {
	c.isVerbose = true
	return c
}

func (c *Cmd) ParseFlags(f Flag) *Cmd {
	c.flags = f

	if f.Args.Profile != "" {
		c.Profile = Cfg().GetProfile(f.Args.Profile)
	}

	if f.Args.Start != "" {
		c.Start = f.Args.Start
	}

	if f.Args.End != "" {
		c.End = f.Args.End
	}

	if f.Args.Output != "" {
		c.Output = file.New(f.Args.Output)
	}

	if f.Args.Input != "" {
		c.Input = file.New(f.Args.Input)
		c.Media = media.NewMedia(f.Args.Input)
	}

	if f.Args.Cover != "" {
		c.Cover = file.New(f.Args.Cover)
		c.Media.AddFile("cover", f.Args.Cover)
	}

	if f.Args.Meta != "" {
		c.Meta = file.New(f.Args.Meta)
		c.Media.SetFFmeta(f.Args.Meta)
	}

	if f.Args.Cue != "" {
		c.Cue = file.New(f.Args.Cue)
		c.Media.SetCue(f.Args.Cue)
	}

	if f.Args.Input != "" {
		c.Media.SetMeta()
	}

	c.PadOutput = Cfg().Defaults.HasPadding()
	c.Padding = Cfg().Defaults.Padding
	c.Num = 1

	if f.Args.ChapNo != 0 {
		c.ChapNo = f.Args.ChapNo
	}
	return c
}

func (c *Cmd) FFmpeg() *ffmpeg.Cmd {
	ffcmd := c.Profile.FFmpegCmd()

	if c.flag.Verbose {
		ffcmd.LogLevel("info")
	}

	if c.flag.Overwrite {
		ffcmd.AppendPreInput("y")
	}

	if c.HasStart() {
		ffcmd.AppendPreInput("ss", c.Start)
	}

	if c.HasEnd() {
		ffcmd.AppendPreInput("to", c.End)
	}

	if c.Media != nil {
		ffcmd.Input(c.Input.Abs)
	}

	if c.HasMeta() {
		ffcmd.FFmeta(c.Meta.Abs)
	}

	if !c.HasOutput() {
		c.Output = file.New(c.Input.Abs)
	}
	ffcmd.Output(c.Output.Abs)

	return ffcmd
}

func (cmd *Cmd) MkTmp() *os.File {
	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	cmd.tmp = f
	return f
}

func (c Cmd) RunBatch() []byte {
	if c.tmp != nil {
		defer os.Remove(c.tmp.Name())
	}

	var buf bytes.Buffer
	for _, cmd := range c.Batch {
		fmt.Println(cmd.String())
		out, err := cmd.Run()
		if err != nil {
			fmt.Printf("cmd string: %s\n", cmd.String())
			log.Fatal(err)
		}

		fmt.Println(cmd.String())

		_, err = buf.Write(out)
		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.Bytes()
}

func (c Cmd) Run() ([]byte, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	cmd, err := c.Build()
	if err != nil {
		return stderr.Bytes(), fmt.Errorf("Cmd failed to build: %v\n", cmd.String())
	}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		return stderr.Bytes(), fmt.Errorf("%v\n", stderr.String())
	}

	return stdout.Bytes(), nil
}
