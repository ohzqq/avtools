package cmd

import (
	"fmt"

	"github.com/ohzqq/fftools/fftools"

	"github.com/integrii/flaggy"
)

type Cli struct{
	Cmd *clir.Cli
	Flags Flags
}

func NewCli() *Cli {
	cli := Cli{
		//Cmd: clir.NewCli("fftools", "FFmpeg tools", "v0.0.1"),
		Flags: initFlags(),
	}
	cli.setFlags()
	//fmt.Printf("%v", cli)
	//cli.Cmd.Action(cli.defaultAction())
	return &cli
}

func (c *Cli) defaultAction() clir.Action {
	return func() error {
		//cmd := fftools.NewCmd().Args(c.Args())
		//c.Cmd.j c 
		//c.Cmd.PrintHelp()
		fmt.Printf("%v", c.Flags["Input"])
		return nil
	}
}

func (c *Cli) Args() fftools.CmdArgs {
	args := make(fftools.CmdArgs)
	for arg, flag := range c.Flags {
		args[arg] = flag.Value
	}
	return args
}

type Flags map[string]*Flag

type Flag struct{
	Name string
	Value string
}

var cmdFlags = map[string]string{
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

func (c *Cli) setFlags() {
	for arg, flag := range cmdFlags {
		//f := c.Flags[arg].Value
		c.Cmd.StringFlag(flag, arg, &c.Flags[arg].Value)
	}
}

func initFlags() Flags {
	flags := make(Flags)
	for arg, _ := range cmdFlags {
		f := Flag{}
		f.Name = cmdFlags[arg]
		flags[arg] = &f
	}
	return flags
}
