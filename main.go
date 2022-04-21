package main

import (
	"fmt"
	//"strings"

	//"github.com/ohzqq/fftools/cmd"
	"github.com/ohzqq/fftools/fftools"

	"github.com/integrii/flaggy"
)

var posInput string
var ext string

func main() {
	fftools.InitConfig()
	fftools.FFcfg()

	var (
		input []string
		output string
		cueSheet string
		cover string
		meta string
		profile = fftools.Cfg.Defaults.Profile
	)

	// Flags
	flaggy.StringSlice(&input, "i", "input", "input")
	flaggy.String(&output, "o", "output", "Set output")
	flaggy.String(&cueSheet, "c", "cue", "set cue sheet")
	flaggy.String(&cover, "C", "cover", "set cover")
	flaggy.String(&meta, "m", "meta", "set ffmetadata")
	flaggy.String(&profile, "p", "profile", "designate profile")

	// Subcommands
	join := joinCmd()
	flaggy.AttachSubcommand(join, 1)

	split := cmdWithInput("split", posInput)
	flaggy.AttachSubcommand(split, 1)

	convert := cmdWithInput("convert", posInput)
	flaggy.AttachSubcommand(convert, 1)

	cueS := cmdWithInput("cue", posInput)
	flaggy.AttachSubcommand(cueS, 1)

	rm := newParentCmd("rm")
	flaggy.AttachSubcommand(rm, 1)
	rmChaps := newChildCmd("chaps", posInput)
	rm.AttachSubcommand(rmChaps, 1)
	rmCover := newChildCmd("cover", posInput)
	rm.AttachSubcommand(rmCover, 1)
	rmMeta := newChildCmd("meta", posInput)
	rm.AttachSubcommand(rmMeta, 1)

	embed := newParentCmd("embed")
	flaggy.AttachSubcommand(embed, 1)
	embedChaps := newChildCmd("chaps", posInput)
	embed.AttachSubcommand(embedChaps, 1)
	embedCover := newChildCmd("cover", posInput)
	embed.AttachSubcommand(embedCover, 1)
	embedMeta := newChildCmd("meta", posInput)
	embed.AttachSubcommand(embedMeta, 1)

	extract := newParentCmd("extract")
	flaggy.AttachSubcommand(extract, 1)
	extractChaps := newChildCmd("chaps", posInput)
	extract.AttachSubcommand(extractChaps, 1)
	extractCover := newChildCmd("cover", posInput)
	extract.AttachSubcommand(extractCover, 1)
	extractMeta := newChildCmd("meta", posInput)
	extract.AttachSubcommand(extractMeta, 1)

	flaggy.Parse()

	// Setup command
	cmd := fftools.NewCmd().Profile(profile)
	//fftools.AllJsonMeta()

	// Input
	if posInput != "" {
		cmd.In(posInput)
	}
	//for _, in := range input {
	//  cmd.In(in)
	//}


	// Handle flags
	if output != "" {
		cmd.Args().Out(output)
	}

	if cueSheet != "" {
		cmd.Args().Cue(cueSheet)
	}

	if cover != "" {
		cmd.Args().Cover(cover)
	}

	if meta != "" {
		cmd.Args().Meta(meta)
	}

	if convert.Used {
		fmt.Println("split")
	}

	if cueS.Used {
		//c := fftools.NewFFProbeCmd().AllJsonMeta(posInput)
		c := fftools.ReadCueSheet(posInput)
		c.Chapters.Timestamps()
		//c.Timestamps()
	}

	if join.Used {
		cmd.Join(ext).Run()
	}

	if split.Used {
		fmt.Println("split")
	}

	if embedChaps.Used {
		fmt.Println("rm cover")
	}

	if embedCover.Used {
		fmt.Println("rm cover")
	}

	if embedMeta.Used{
		fmt.Println("rm cover")
	}

	if extractChaps.Used {
		fmt.Println("rm cover")
	}

	if extractCover.Used {
		fmt.Println("rm cover")
	}

	if extractMeta.Used{
		fmt.Println("rm cover")
	}

	if rmChaps.Used {
		fmt.Println("rm cover")
	}

	if rmCover.Used {
		fmt.Println("rm cover")
	}

	if rmMeta.Used{
		fmt.Println("rm cover")
	}

	//fmt.Println(cmd.Cmd().String())
}

func joinCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("join")
	cmd.AddPositionalValue(&ext, "extension", 1, true, "specify the extension of the files")
	return cmd
}

func cmdWithInput(name string, input string) *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand(name)
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func newParentCmd(name string) *flaggy.Subcommand {
	return flaggy.NewSubcommand(name)
}

func newChildCmd(name string, input string) *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand(name)
	cmd.AddPositionalValue(&input, "input", 1, true, "input for the command")
	return cmd
}
