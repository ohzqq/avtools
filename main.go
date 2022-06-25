package main

import (
	"fmt"
	//"path/filepath"
	"log"
	//"strings"

	"github.com/ohzqq/avtools/avtools"

	"github.com/integrii/flaggy"
)


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	avtools.InitConfig()

	cmd := avtools.NewCmd()

	// Flags
	flaggy.StringSlice(&cmd.InputSlice, "i", "input", "input")
	flaggy.String(&cmd.Output, "o", "output", "Set output")
	flaggy.Bool(&cmd.Overwrite, "y", "yes", "overwrite")
	flaggy.Bool(&cmd.Verbose, "v", "", "print command")
	flaggy.String(&cmd.Profile, "p", "profile", "designate profile")

	// Subcommands
	test := flaggy.NewSubcommand("test")
	test.AddPositionalValue(&cmd.Input, "file", 1, true, "specify the extension of the files")
	flaggy.AttachSubcommand(test, 1)

	join := flaggy.NewSubcommand("join")
	join.AddPositionalValue(&cmd.Ext, "extension", 1, true, "specify the extension of the files")
	flaggy.AttachSubcommand(join, 1)

	cut := flaggy.NewSubcommand("cut")
	cut.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	cut.String(&cmd.MetaFile, "m", "meta", "use ffmetadata file as markers")
	cut.String(&cmd.CueFile, "c", "cue", "use cue sheet as markers")
	cut.String(&cmd.Start, "ss", "start", "the start time")
	cut.String(&cmd.End, "to", "end", "the end time")
	cut.Int(&cmd.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(cut, 1)

	split := flaggy.NewSubcommand("split")
	split.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	split.String(&cmd.CueFile, "c", "cue", "use cue sheet as markers")
	split.String(&cmd.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(split, 1)

	convert := flaggy.NewSubcommand("convert")
	convert.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	flaggy.AttachSubcommand(convert, 1)

	cueS := flaggy.NewSubcommand("cue")
	cueS.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	flaggy.AttachSubcommand(cueS, 1)

	rm := flaggy.NewSubcommand("rm")
	rm.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	rm.Bool(&cmd.MetaFlag, "m", "meta", "Remove meta")
	rm.Bool(&cmd.ChapFlag, "c", "chaps", "Remove chapters")
	rm.Bool(&cmd.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(rm, 1)

	show := flaggy.NewSubcommand("show")
	show.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	show.Bool(&cmd.MetaFlag, "m", "meta", "Remove meta")
	show.Bool(&cmd.ChapFlag, "c", "chaps", "Remove chapters")
	show.Bool(&cmd.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(show, 1)

	extract := flaggy.NewSubcommand("extract")
	extract.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	extract.Bool(&cmd.MetaFlag, "m", "meta", "extract meta")
	extract.Bool(&cmd.ChapFlag, "c", "chaps", "extract chapters")
	extract.Bool(&cmd.CoverFlag, "a", "art", "extract album art")
	flaggy.AttachSubcommand(extract, 1)

	update := flaggy.NewSubcommand("update")
	update.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	update.String(&cmd.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	update.String(&cmd.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(update, 1)

	flaggy.Parse()

	// Input
	if cmd.Input != "" {
		cmd.Media.Input(cmd.Input)
	}

	// Handle flags
	if cmd.Overwrite {
		cmd.Media.Overwrite = true
	}

	if cmd.CueFile != "" {
		cmd.Media.SetChapters(avtools.ReadCueSheet(cmd.CueFile))
	}

	if test.Used {
		fmt.Printf("%+V\n", avtools.Cfg())
		//cmd := avtools.RmChapters(media)
		//cmd.Run()
	}

	switch {
	case join.Used:
		cmd.Join()
	case show.Used:
		cmd.Show()
	case cut.Used:
		cmd.Cut(cmd.Start, cmd.End, cmd.ChapNo).Run()
	case split.Used:
		err := cmd.Split()
		if err != nil {
			log.Fatal(err)
		}
	case rm.Used:
		cmd.Remove()
	case extract.Used:
		cmd.Extract()
	case update.Used:
		cmd.Update()
	case convert.Used:
		fmt.Println("split")
	case cueS.Used:
		c := avtools.ReadCueSheet(cmd.Input)
		fmt.Printf("%V", c)
	}
}
