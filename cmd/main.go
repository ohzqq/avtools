package cmd

import (
	"fmt"

	"github.com/ohzqq/fftools/fftools"

	"github.com/leaanthony/clir"
)

type Cli struct{
	Cmd *clir.Cli
	Flags flags
}

const (
	Pre = "pre"
	Input = "i"
	Post = "post"
	VideoCodec = "vc"
	VideoParams = "vp"
	VideoFilters = "vf"
	AudioCodec = "ac"
	AudioParams = "ap"
	AudioFilters = "af"
	FilterCompex = "filter"
	Output = "o"
)

type flags map[string]Flag

type Flag struct{
	Changed bool
	Flag string
	Value string
	Type string
}

var Flags = map[string]string{
	"Pre": "pre",
	"Input": "i",
	"Post": "post",
	"VideoCodec": "vc",
	"VideoParams": "vp",
	"VideoFilters": "vf",
	"AudioCodec": "ac",
	"AudioParams": "ap",
	"AudioFilters": "af",
	"FilterCompex": "filter",
	"Output": "o",
}

var CmdFlags = make(flags)

func (f Flag) HasChanged() bool { return f.Changed }
func (f Flag) Name() string { return f.Flag }
func (f Flag) ValueString() string { return f.Value }
func (f Flag) ValueType() string { return f.Type }

func NewCli() *Cli {
	cli := Cli{
		Cmd: clir.NewCli("fftools", "FFmpeg tools", "v0.0.1"),
		Flags: initFlags(),
	}
	cli.setFlags()
	fmt.Printf("%v", cli)
	cli.Cmd.Action(cli.defaultAction())
	return &cli
}

func (c *Cli) defaultAction() clir.Action {
	return func() error {
		c.Cmd.PrintHelp()
		//fmt.Println("Hello world")
		return nil
	}
}

func (c *Cli) setFlags() {
	for arg, flag := range Flags {
		f := c.Flags[arg].Value
		c.Cmd.StringFlag(flag, arg, &f)
	}
}

func initFlags() flags {
	f := make(flags)
	for _, arg := range fftools.ArgOrder {
		flag := Flag{}
		flag.Flag = Flags[arg]
		f[arg] = flag
	}
	return f
}
