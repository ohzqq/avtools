package cmd

import (
	"fmt"

	//"github.com/ohzqq/fftools/fftools"

	"github.com/leaanthony/clir"
)

type Cli struct{
	Cmd *clir.Cli
	Flags flags
}

type flags map[string]*Flag

type Flag struct{
	Name string
	Value string
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
	"Cue": "cue",
	"Cover": "c",
	"FFmetadata": "m",
	"Verbosity": "v",
	"Profile": "p",
}

func NewCli() *Cli {
	cli := Cli{
		Cmd: clir.NewCli("fftools", "FFmpeg tools", "v0.0.1"),
		Flags: initFlags(),
	}
	cli.setFlags()
	//fmt.Printf("%v", cli)
	cli.Cmd.Action(cli.defaultAction())
	return &cli
}

func (c *Cli) defaultAction() clir.Action {
	return func() error {
		//c.Cmd.j c 
		//c.Cmd.PrintHelp()
		fmt.Printf("%v", c.Flags["Cover"])
		return nil
	}
}

func (c *Cli) setFlags() {
	for arg, flag := range Flags {
		c.Cmd.StringFlag(flag, arg, &c.Flags[arg].Value)
	}
}

func initFlags() flags {
	flags := make(flags)
	for arg, _ := range Flags {
		f := Flag{}
		f.Name = Flags[arg]
		flags[arg] = &f
	}
	return flags
}
