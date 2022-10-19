package tool

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/media"
)

type Cmd struct {
	Flag
	Media     *media.Media
	isVerbose bool
	cwd       string
	Batch     []Command
	bin       string
	args      []string
	tmpFile   string
	num       int
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
	c.Media = media.NewMedia(i)
	return c
}

func (c Cmd) String() string {
	return strings.Join(c.args, " ")
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

func (c *Cmd) SetFlags(f Flag) *Cmd {
	c.Flag = f
	c.Media = f.Media()
	return c
}

func (c *Cmd) FFmpeg() *ffmpeg.Cmd {
	cmd := Cfg().GetProfile("default").FFmpegCmd()

	if c.Flag.Args.HasProfile() {
		cmd = Cfg().GetProfile(c.Flag.Args.Profile).FFmpegCmd()
	}

	if c.Bool.Verbose {
		cmd.LogLevel("info")
	}

	if c.Bool.Overwrite {
		cmd.AppendPreInput("y")
	}

	if c.Args.HasStart() {
		cmd.AppendPreInput("ss", c.Args.Start)
	}

	if c.Args.HasEnd() {
		cmd.AppendPreInput("to", c.Args.End)
	}

	if c.Media != nil {
		cmd.Input(c.Media.Input)
	}

	if c.Args.HasMeta() {
		cmd.FFmeta(c.Args.Meta)
	}

	var out *Output
	if c.Args.HasOutput() {
		out = NewOutput(c.Args.Output)
	} else {
		out = NewOutput(c.Args.Input)
	}
	cmd.Output(out.String())

	return cmd
}

func (cmd *Cmd) Tmp(f string) *Cmd {
	cmd.tmpFile = f
	return cmd
}

func (c Cmd) RunBatch() []byte {
	var buf bytes.Buffer
	for _, cmd := range c.Batch {
		out, err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

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
