package ff

import (
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Cmd struct {
	Filters Filters `yaml:"filters"`
	Output
	Input
}

func New(profile ...string) Cmd {
	if len(profile) > 0 {
		pro := GetProfile(profile[0])
		return pro
	}

	def := Cmd{
		Filters: make(Filters),
		Output:  NewOutput(),
		Input: NewInput(ffmpeg.KwArgs{
			"loglevel":    "error",
			"hide_banner": "",
		}),
	}
	return def
}

func (cmd *Cmd) In(file string, args ...ffmpeg.KwArgs) *Cmd {
	kwargs := []ffmpeg.KwArgs{cmd.Input.Args}
	kwargs = append(kwargs, args...)
	cmd.Input.Args = ffmpeg.MergeKwArgs(kwargs)
	cmd.File = file
	return cmd
}

func (cmd *Cmd) Compile() *exec.Cmd {
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

	ffArgs := output.GetArgs()

	var args []string

	args = append(args, ffArgs[:inArgs]...)

	if meta, ok := cmd.Input.Args["meta"]; ok {
		args = append(args, "-i", meta.(string))
	}

	args = append(args, ffArgs[inArgs:outArgs]...)

	if label, ok := cmd.Input.Args["map_metadata"]; ok {
		args = append(args, "-map_metadata", label.(string))
	}

	args = append(args, ffArgs[outArgs:]...)

	exe := exec.Command("ffmpeg", args...)

	return exe
}

func (cmd Cmd) Run() error {
	err := cmd.Compile().Run()
	if err != nil {
		return err
	}
	return nil
}
