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
	test := cmdWithInput("test")
	flaggy.AttachSubcommand(test, 1)

	join := joinCmd()
	flaggy.AttachSubcommand(join, 1)

	split := cmdWithInput("split")
	flaggy.AttachSubcommand(split, 1)

	convert := cmdWithInput("convert")
	flaggy.AttachSubcommand(convert, 1)

	cueS := cmdWithInput("cue")
	flaggy.AttachSubcommand(cueS, 1)

	rm := newParentCmd("rm")
	flaggy.AttachSubcommand(rm, 1)
	rmChaps := newChildCmd("chaps")
	rm.AttachSubcommand(rmChaps, 1)
	rmCover := newChildCmd("cover")
	rm.AttachSubcommand(rmCover, 1)
	rmMeta := newChildCmd("meta")
	rm.AttachSubcommand(rmMeta, 1)

	embed := newParentCmd("embed")
	flaggy.AttachSubcommand(embed, 1)
	embedChaps := newChildCmd("chaps")
	embed.AttachSubcommand(embedChaps, 1)
	embedCover := newChildCmd("cover")
	embed.AttachSubcommand(embedCover, 1)
	embedMeta := newChildCmd("meta")
	embed.AttachSubcommand(embedMeta, 1)

	extract := newParentCmd("extract")
	flaggy.AttachSubcommand(extract, 1)
	extractChaps := newChildCmd("chaps")
	extract.AttachSubcommand(extractChaps, 1)
	extractCover := newChildCmd("cover")
	extract.AttachSubcommand(extractCover, 1)
	extractMeta := newChildCmd("meta")
	extract.AttachSubcommand(extractMeta, 1)

	flaggy.Parse()

	// Setup command
	cmd := fftools.NewCmd().Profile(profile)

	// Input
	if posInput != "" {
		cmd.In(fftools.NewMedia(posInput))
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

	if test.Used {
		//file := fftools.NewMedia(posInput).WithMeta()
		//fmt.Printf("%v\n", file.HasChapters())
		//file.ReadMeta()
		ch := cmd.GetChapters()
		fmt.Printf("%T\n", ch)
		//fmt.Printf("%V\n", file.Meta.Tags.Title)
		//fmt.Printf("%v\n", file.HasChapters())
		//fmt.Printf("%v\n", )
	}

	if convert.Used {
		fmt.Println("split")
	}

	if cueS.Used {
		//c := fftools.AllJsonMeta(posInput)
		c := fftools.ReadCueSheet(posInput)
		//c := cmd.Meta()
		fmt.Printf("%V", c.Chapters)
		//c.Chapters.Timestamps()
		//c.Timestamps()
	}

	if join.Used {
		fftools.Join(ext).Profile(profile).Run()
	}

	if split.Used {
		cmd.Split()
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
		fftools.WriteFFmetadata(posInput)
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

func cmdWithInput(name string) *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand(name)
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func newParentCmd(name string) *flaggy.Subcommand {
	return flaggy.NewSubcommand(name)
}

func newChildCmd(name string) *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand(name)
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}
