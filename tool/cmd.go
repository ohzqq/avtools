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
	output    file.File
	Input     file.File
	Cover     file.File
	Meta      file.File
	Cue       file.File
	Media     *media.Media
	Num       int
	PadOutput bool
	Padding   string
	ChapNo    int
	flag      Flag
	output    Output
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
		cwd: cwd,
	}
}

func (c *Cmd) Input(i string) *Cmd {
	c.Args.Input = file.New(i)
	c.Args.Media = media.NewMedia(i)
	return c
}

func (c Cmd) String() string {
	return strings.Join(c.args, " ")
}

func (c Cmd) Media() *media.Media {
	return c.Args.Media
}

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
	c.flag = f

	if f.Args.Profile != "" {
		c.Profile = f.Args.GetProfile()
	}

	if f.Args.Start != "" {
		c.Start = f.Args.Start
	}
	if f.Args.End != "" {
		c.End = f.Args.End
	}
	if f.Args.Output != "" {
		c.output = file.New(f.Args.Output)
	}
	if f.Args.Input != "" {
		c.Input = file.New(f.Args.Input)
	}
	if f.Args.Cover != "" {
		c.Cover = file.New(f.Args.Cover)
	}
	if f.Args.Meta != "" {
		c.Meta = file.New(f.Args.Meta)
	}
	if f.Args.Cue != "" {
		c.Cue = file.New(f.Args.Cue)
	}
	c.Media = f.Media()
	c.PadOutput = Cfg().Defaults.HasPadding()
	c.Padding = Cfg().Defaults.Padding
	c.Num = 1
	c.ChapNo = f.Args.ChapNo
	return c
}

func (c *Cmd) FFmpeg() *ffmpeg.Cmd {
	ffcmd := Cfg().GetProfile("default").FFmpegCmd()

	if c.flag.Args.HasProfile() {
		ffcmd = Cfg().GetProfile(c.flag.Args.Profile).FFmpegCmd()
	}

	if c.flag.Bool.Verbose {
		ffcmd.LogLevel("info")
	}

	if c.flag.Bool.Overwrite {
		ffcmd.AppendPreInput("y")
	}

	if c.flag.Args.HasStart() {
		ffcmd.AppendPreInput("ss", c.Args.Start)
	}

	if c.flag.Args.HasEnd() {
		ffcmd.AppendPreInput("to", c.Args.End)
	}

	if c.Args.Media != nil {
		ffcmd.Input(c.Args.Input.Abs)
	}

	if c.flag.Args.HasMeta() {
		ffcmd.FFmeta(c.flag.Args.Meta)
	}

	if !c.flag.Args.HasOutput() {
		c.flag.Args.Output = c.flag.Args.Input
	}
	out := NewOutput(c.flag.Args.Output)
	ffcmd.Output(out.String())

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
