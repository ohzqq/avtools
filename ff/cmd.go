package ff

import (
	"fmt"

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

func (cmd *Cmd) Compile() *ffmpeg.Stream {
	input := NewInput(cmd.Input.Args)
	in := input.Compile(cmd.File)

	for _, filter := range cmd.Filters.Compile() {
		fmt.Printf("filter %+V\n", filter)
		in = filter(in)
	}

	output := cmd.Output.Compile(in)

	return output
}

func (cmd Cmd) Run() error {
	err := cmd.Compile().ErrorToStdOut().Run()
	if err != nil {
		return err
	}
	return nil
}
