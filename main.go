package main

import (
	"fmt"
	//"strings"
	//"log"

	//"github.com/ohzqq/fftools/cmd"
	"github.com/ohzqq/fftools/fftools"

	"github.com/integrii/flaggy"
)

func main() {
	fftools.InitConfig()
	fftools.FFcfg()

	// Flags
	var input []string
	flaggy.StringSlice(&input, "i", "input", "input")

	var output string
	flaggy.String(&output, "o", "output", "Set output")

	var cue string
	flaggy.String(&cue, "c", "cue", "set cue sheet")

	var cover string
	flaggy.String(&cover, "C", "cover", "set cover")

	var meta string
	flaggy.String(&meta, "m", "meta", "set ffmetadata")

	var profile = fftools.Cfg.Defaults.Profile
	flaggy.String(&profile, "p", "profile", "designate profile")

	flaggy.Parse()

	cmd := fftools.NewCmd().Profile(profile)
	for _, in := range input {
		cmd.In(in)
	}

	if output != "" {
		cmd.Args().Out(output)
	}

	if cue != "" {
		cmd.Args().Cue(cue)
	}

	if cover != "" {
		cmd.Args().Cover(cover)
	}

	if meta != "" {
		cmd.Args().Meta(meta)
	}

	cmd.Args() //.VCodec("libx264")
	//c := fftools.NewCmd().Input(input)

	fmt.Printf("%v", cmd.Args())
	fmt.Printf("%v", cmd.Cmd())
}
