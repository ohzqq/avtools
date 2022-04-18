package cmd

import (
	//"fmt"

	"github.com/ohzqq/fftools/fftools"

	"github.com/integrii/flaggy"
)

type Cli struct{
	Parser *flaggy.Parser
	Flags Flags
}

func NewCli() *Cli {
	cli := Cli{
		Parser: flaggy.NewParser("fftools"),
		Flags: initFlags(),
	}
	cli.setFlags()
	//fmt.Printf("%v", cli.Flags["af"])
	//cli.Cmd.Action(cli.defaultAction())
	return &cli
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

func (c Cli) setFlags() {
	for arg, flag := range cmdFlags {
		//f := c.Flags[arg].Value
		flaggy.String(&c.Flags[arg].Value, flag, arg, arg)
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
