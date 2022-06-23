package main

import (
	"fmt"
	//"path/filepath"
	//"log"
	//"strings"

	"github.com/ohzqq/avtools/avtools"

	"github.com/integrii/flaggy"
)

var posInput string
var ext string

func main() {
	avtools.InitConfig()

	var (
		input []string
		output string
		cueSheet string
		cover string
		overwrite bool
		meta string
		profile = avtools.Cfg().DefaultProfile()
	)

	// Flags
	flaggy.StringSlice(&input, "i", "input", "input")
	flaggy.String(&output, "o", "output", "Set output")
	flaggy.Bool(&overwrite, "y", "yes", "overwrite")
	//flaggy.String(&cueSheet, "c", "cue", "set cue sheet")
	//flaggy.String(&cover, "C", "cover", "set cover")
	//flaggy.String(&meta, "m", "meta", "set ffmetadata")
	flaggy.String(&profile, "p", "profile", "designate profile")

	// Subcommands
	test := cmdWithInput("test")
	flaggy.AttachSubcommand(test, 1)

	join := joinCmd()
	flaggy.AttachSubcommand(join, 1)

	split := cmdWithInput("split")
	split.String(&cueSheet, "c", "cue", "use cue sheet as markers")
	flaggy.AttachSubcommand(split, 1)

	convert := cmdWithInput("convert")
	flaggy.AttachSubcommand(convert, 1)

	cueS := cmdWithInput("cue")
	flaggy.AttachSubcommand(cueS, 1)

	var (
		xMeta bool
		xChaps bool
		xCover bool
	)
	rm := newChildCmd("rm")
	rm.Bool(&xMeta, "m", "meta", "Remove meta")
	rm.Bool(&xChaps, "c", "chaps", "Remove chapters")
	rm.Bool(&xCover, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(rm, 1)

	extract := newChildCmd("extract")
	extract.Bool(&xMeta, "m", "meta", "extract meta")
	extract.Bool(&xChaps, "c", "chaps", "extract chapters")
	extract.Bool(&xCover, "a", "art", "extract album art")
	flaggy.AttachSubcommand(extract, 1)

	update := newChildCmd("update")
	update.String(&meta, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	update.String(&cover, "a", "art", "update album art")
	flaggy.AttachSubcommand(update, 1)

	flaggy.Parse()

	// Setup command
	//cmd := avtools.NewCmd().Profile(profile)

	// Input
	var media *avtools.Media
	if posInput != "" {
		media = avtools.NewMedia(posInput).WithMeta()
	}
	//for _, in := range input {
	//  cmd.In(in)
	//}


	// Handle flags
	if overwrite {
		//cmd.Args().OverWrite()
		media.Overwrite = true
	}

	if output != "" {
		//cmd.Args().Out(output)
	}

	if cueSheet != "" {
		//media.Cue = filepath.Base(cueSheet)
		//cmd.Args().Cue(cueSheet)
		media.SetChapters(avtools.ReadCueSheet(cueSheet))
	}

	//if cover != "" {
		//media.Cover = filepath.Base(cover)
		//cmd.Args().Cover(cover)
	//}

	//if meta != "" {
		//media.SetMeta(avtools.ReadFFmetadata(meta))
		//cmd.Args().Meta(meta)
	//}

	//cmd.In(media)

	if test.Used {
		fmt.Printf("%+V\n", avtools.Cfg())
		//cmd := avtools.RmChapters(media)
		//cmd.Run()
		//fmt.Printf("%V\n", media.HasStreams())
		//fmt.Printf("%V\n", file.Meta.Tags.Title)
		//fmt.Printf("%V\n", media.Meta.Chapters)
		//fmt.Printf("%v\n", cmd.String())
	}

	if convert.Used {
		fmt.Println("split")
	}

	if cueS.Used {
		//c := avtools.AllJsonMeta(posInput)
		c := avtools.ReadCueSheet(posInput)
		//c := cmd.Meta()
		fmt.Printf("%V", c)
		//c.Chapters.Timestamps()
		//c.Timestamps()
	}

	if join.Used {
		avtools.Join(ext).Profile(profile).Run()
	}

	if split.Used {
		media.Split(cueSheet)
	}

	if update.Used {
		media.Update(cover, meta)
		//fmt.Printf("%+V\n", cmd.String())
		//cmd.Run()
	}

	if extract.Used {
		media.Extract(xChaps, xCover, xMeta)
	}

	if rm.Used {
		media.Remove(xChaps, xCover, xMeta).Run()
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
