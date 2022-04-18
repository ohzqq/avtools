package main

import (
	"fmt"
	//"log"

	//"github.com/ohzqq/fftools/cmd"
	"github.com/ohzqq/fftools/fftools"

	"github.com/integrii/flaggy"
)

func main() {
	fftools.InitConfig()
	fftools.FFcfg()
	//c := fftools.NewCmd().Args(fftools.Cfg.Profiles["convert"])
	//fmt.Printf("%v", c.Cmd())
	//cli := cmd.NewCli()
	//fmt.Printf("%v", fftools.Cfg.Profiles)
	//cli.Parser.ShowHelp()
	var input []string
	flaggy.StringSlice(&input, "i", "input", "input")
	cli := Cli{}
	cli.Flags = initFlags()
	for arg, flag := range cmdFlags {
		//f := c.Flags[arg].Value
		flaggy.String(&cli.Flags[arg].Value, flag, arg, arg)
	}
	flaggy.Parse()
	fmt.Printf("%v", cli.Args())
}

type Cli struct{
	Parser *flaggy.Parser
	Flags Flags
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
	//"Input": "i",
	"Post": "post",
	"VideoCodec": "c:v",
	"VideoParams": "vp",
	"VideoFilters": "vf",
	"AudioCodec": "c:a",
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

