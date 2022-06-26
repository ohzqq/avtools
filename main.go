package main

import (
	//"fmt"
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

	// Join command
	join := flaggy.NewSubcommand("join")
	join.AddPositionalValue(&cmd.Ext, "extension", 1, true, "specify the extension of the files")
	flaggy.AttachSubcommand(join, 1)
	// Join alias
	j := flaggy.NewSubcommand("j")
	j.AddPositionalValue(&cmd.Ext, "extension", 1, true, "specify the extension of the files")
	flaggy.AttachSubcommand(j, 1)

	// Cut command
	cut := flaggy.NewSubcommand("cut")
	cut.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	cut.String(&cmd.MetaFile, "m", "meta", "use ffmetadata file as markers")
	cut.String(&cmd.CueFile, "c", "cue", "use cue sheet as markers")
	cut.String(&cmd.Start, "ss", "start", "the start time")
	cut.String(&cmd.End, "to", "end", "the end time")
	cut.Int(&cmd.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(cut, 1)
	// Cut alias
	c := flaggy.NewSubcommand("c")
	c.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	c.String(&cmd.MetaFile, "m", "meta", "use ffmetadata file as markers")
	c.String(&cmd.CueFile, "c", "cue", "use cue sheet as markers")
	c.String(&cmd.Start, "ss", "start", "the start time")
	c.String(&cmd.End, "to", "end", "the end time")
	c.Int(&cmd.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(c, 1)

	// Split command
	split := flaggy.NewSubcommand("split")
	split.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	split.String(&cmd.CueFile, "c", "cue", "use cue sheet as markers")
	split.String(&cmd.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(split, 1)
	// Split alias
	s := flaggy.NewSubcommand("s")
	s.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	s.String(&cmd.CueFile, "c", "cue", "use cue sheet as markers")
	s.String(&cmd.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(s, 1)

	// Convert command
	//convert := flaggy.NewSubcommand("convert")
	//convert.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	//flaggy.AttachSubcommand(convert, 1)

	// Remove command
	remove := flaggy.NewSubcommand("remove")
	remove.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	remove.Bool(&cmd.MetaFlag, "m", "meta", "Remove meta")
	remove.Bool(&cmd.ChapFlag, "c", "chaps", "Remove chapters")
	remove.Bool(&cmd.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(remove, 1)
	// Remove alias
	rm := flaggy.NewSubcommand("rm")
	rm.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	rm.Bool(&cmd.MetaFlag, "m", "meta", "Remove meta")
	rm.Bool(&cmd.ChapFlag, "c", "chaps", "Remove chapters")
	rm.Bool(&cmd.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(rm, 1)

	// Show command
	show := flaggy.NewSubcommand("show")
	show.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	show.Bool(&cmd.MetaFlag, "m", "meta", "Remove meta")
	show.Bool(&cmd.ChapFlag, "c", "chaps", "Remove chapters")
	show.Bool(&cmd.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(show, 1)

	// Extract command
	extract := flaggy.NewSubcommand("extract")
	extract.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	extract.Bool(&cmd.MetaFlag, "m", "meta", "extract meta")
	extract.Bool(&cmd.ChapFlag, "c", "chaps", "extract chapters")
	extract.Bool(&cmd.CoverFlag, "a", "art", "extract album art")
	flaggy.AttachSubcommand(extract, 1)
	// Extract alias
	x := flaggy.NewSubcommand("x")
	x.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	x.Bool(&cmd.MetaFlag, "m", "meta", "extract meta")
	x.Bool(&cmd.ChapFlag, "c", "chaps", "extract chapters")
	x.Bool(&cmd.CoverFlag, "a", "art", "extract album art")
	flaggy.AttachSubcommand(x, 1)

	// Update command
	update := flaggy.NewSubcommand("update")
	update.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	update.String(&cmd.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	update.String(&cmd.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(update, 1)
	// Update alias
	u := flaggy.NewSubcommand("u")
	u.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	u.String(&cmd.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	u.String(&cmd.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(u, 1)

	flaggy.Parse()

	// Input
	if cmd.Input != "" {
		cmd.Media.Input(cmd.Input)
	}

	// Handle flags
	if cmd.Overwrite {
		cmd.Media.Overwrite = true
	}

	switch {
	case join.Used, j.Used:
		cmd.Join()
	case show.Used:
		cmd.Show()
	case cut.Used, c.Used:
		cmd.Cut(cmd.Start, cmd.End, cmd.ChapNo).Run()
	case split.Used, s.Used:
		err := cmd.Split()
		if err != nil {
			log.Fatal(err)
		}
	case rm.Used, remove.Used:
		cmd.Remove()
	case extract.Used, x.Used:
		cmd.Extract()
	case update.Used, u.Used:
		cmd.Update()
	}
}
