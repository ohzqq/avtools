package ff

import (
	"bytes"
	"fmt"
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Cmd struct {
	Filters Filters `yaml:"filters"`
	Output
	Input
	args []string
}

func New(profile ...string) Cmd {
	pro := "default"
	if len(profile) > 0 {
		pro = profile[0]
	}

	return GetProfile(pro)
}

func (cmd *Cmd) In(file string, args ...ffmpeg.KwArgs) *Cmd {
	kwargs := []ffmpeg.KwArgs{cmd.Input.Args}
	kwargs = append(kwargs, args...)
	cmd.Input.Args = ffmpeg.MergeKwArgs(kwargs)
	cmd.File = file
	return cmd
}

func (cmd *Cmd) Compile() *Cmd {
	input := NewInput(cmd.Input.Args)
	in := input.Compile(cmd.File)
	inArgs := len(in.GetArgs())

	for _, filter := range cmd.Filters.Compile() {
		in = filter(in)
	}
	fArgs := len(in.GetArgs()[inArgs:])

	output := cmd.Output.Compile(in)

	outArgs := inArgs
	if fArgs > 0 {
		fArgs = fArgs + 2
		outArgs = outArgs + fArgs
	}

	//output.Compile()

	ffArgs := output.GetArgs()

	cmd.args = append(cmd.args, ffArgs[:inArgs]...)

	if meta, ok := cmd.Input.Args["meta"]; ok {
		cmd.args = append(cmd.args, "-i", meta.(string))
	}

	cmd.args = append(cmd.args, ffArgs[inArgs:outArgs]...)

	if label, ok := cmd.Input.Args["map_metadata"]; ok {
		cmd.args = append(cmd.args, "-map_metadata", label.(string))
	}

	if label, ok := cmd.Input.Args["map_chapters"]; ok {
		cmd.args = append(cmd.args, "-map_chapters", label.(string))
	}

	cmd.args = append(cmd.args, ffArgs[outArgs:]...)

	//fmt.Printf("args %+V\n", args)

	return cmd
}

func (c Cmd) Run() error {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	cmd := exec.Command("ffmpeg", c.args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v\n%v\n", stderr.String(), cmd.String())
	}

	if len(stdout.Bytes()) > 0 {
		fmt.Println(stdout.String())
	}

	return nil
}
