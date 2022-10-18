package media

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ohzqq/avtools/ffmpeg"
)

type Cmd struct {
	Flag
	Media   *Media
	verbose bool
	cwd     string
	exec    *exec.Cmd
	Batch   []*exec.Cmd
	tmpFile string
	num     int
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
	c.Media = NewMedia(i)
	return c
}

func (c *Cmd) Exec(bin string, args []string) *Cmd {
	cmd := exec.Command(bin, args...)
	c.AddCmd(cmd)
	return c
}

func (c *Cmd) AddCmd(cmd *exec.Cmd) *Cmd {
	c.Batch = append(c.Batch, cmd)
	return c
}

func (c *Cmd) Verbose() *Cmd {
	c.verbose = true
	return c
}

func (c *Cmd) SetFlags(f Flag) *Cmd {
	c.Flag = f
	c.Media = f.Media()
	return c
}

func (c *Cmd) NewFFmpegCmd() *ffmpeg.Cmd {
	cmd := ffmpeg.New()

	if v := Cfg().Defaults.LogLevel; v != "" {
		cmd.LogLevel(v)
	}

	if c.Bool.Verbose {
		cmd.LogLevel("info")
	}

	if Cfg().Defaults.Overwrite || c.Bool.Overwrite {
		cmd.AppendPreInput("y")
	}

	if c.Args.HasStart() {
		cmd.AppendPreInput("ss", c.Args.Start)
	}

	if c.Args.HasEnd() {
		cmd.AppendPreInput("to", c.Args.End)
	}

	//if c.Args.HasInput() {
	if c.Media != nil {
		cmd.Input(c.Media.input)
	}

	if c.Args.HasMeta() {
		cmd.FFmeta(c.Args.Meta)
	}

	if c.Args.HasOutput() {
		cmd.Output("tmp")
		//cmd.Output(c.Args.Output)
	} else {
		cmd.Output("tmp")
		//cmd.Output(OutputFromInput(f.Args.Input).String())
	}

	return cmd
}

func (cmd *Cmd) Tmp(f string) *Cmd {
	cmd.tmpFile = f
	return cmd
}

func (cmd Cmd) Run() []byte {
	if cmd.tmpFile != "" {
		defer os.Remove(cmd.tmpFile)
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

	if cmd.verbose {
		fmt.Println(cmd.exec.String())
	}
	return nil
}
