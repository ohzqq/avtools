package main

import (
	"fmt"
	"log"

	"github.com/ohzqq/avtools/avtools"

	"github.com/integrii/flaggy"
)
var _ = fmt.Printf

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	avtools.InitConfig()

	cmd := avtools.NewCmd()

	// Flags
	flaggy.String(&cmd.Flags.Output, "o", "output", "Set output")
	flaggy.Bool(&cmd.Flags.Overwrite, "y", "yes", "overwrite")
	flaggy.Bool(&cmd.Flags.Verbose, "v", "verbose", "print command")
	flaggy.String(&cmd.Flags.Profile, "p", "profile", "designate profile")

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
	cut.String(&cmd.Flags.MetaFile, "m", "meta", "use ffmetadata file as markers")
	cut.String(&cmd.Flags.CueFile, "c", "cue", "use cue sheet as markers")
	cut.String(&cmd.Flags.Start, "ss", "start", "the start time")
	cut.String(&cmd.Flags.End, "to", "end", "the end time")
	cut.Int(&cmd.Flags.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(cut, 1)
	// Cut alias
	c := flaggy.NewSubcommand("c")
	c.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	c.String(&cmd.Flags.MetaFile, "m", "meta", "use ffmetadata file as markers")
	c.String(&cmd.Flags.CueFile, "c", "cue", "use cue sheet as markers")
	c.String(&cmd.Flags.Start, "ss", "start", "the start time")
	c.String(&cmd.Flags.End, "to", "end", "the end time")
	c.Int(&cmd.Flags.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(c, 1)

	// Split command
	split := flaggy.NewSubcommand("split")
	split.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	split.String(&cmd.Flags.CueFile, "c", "cue", "use cue sheet as markers")
	split.String(&cmd.Flags.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(split, 1)
	// Split alias
	s := flaggy.NewSubcommand("s")
	s.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	s.String(&cmd.Flags.CueFile, "c", "cue", "use cue sheet as markers")
	s.String(&cmd.Flags.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(s, 1)

	// Convert command
	//convert := flaggy.NewSubcommand("convert")
	//convert.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	//flaggy.AttachSubcommand(convert, 1)

	// Remove command
	remove := flaggy.NewSubcommand("remove")
	remove.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	remove.Bool(&cmd.Flags.MetaSwitch, "m", "meta", "Remove meta")
	remove.Bool(&cmd.Flags.ChapSwitch, "c", "chaps", "Remove chapters")
	remove.Bool(&cmd.Flags.CoverSwitch, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(remove, 1)
	// Remove alias
	rm := flaggy.NewSubcommand("rm")
	rm.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	rm.Bool(&cmd.Flags.MetaSwitch, "m", "meta", "Remove meta")
	rm.Bool(&cmd.Flags.ChapSwitch, "c", "chaps", "Remove chapters")
	rm.Bool(&cmd.Flags.CoverSwitch, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(rm, 1)

	// Show command
	show := flaggy.NewSubcommand("show")
	show.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	show.AddPositionalValue(&cmd.Action, "action", 2, false, "input for the command")
	show.String(&cmd.Flags.Start, "ss", "start", "the start time")
	show.String(&cmd.Flags.End, "to", "end", "the end time")
	show.Int(&cmd.Flags.ChapNo, "n", "num", "chapter to cut")
	show.String(&cmd.Flags.MetaFile, "m", "metaFile", "use ffmetadata file as markers")
	show.String(&cmd.Flags.CueFile, "c", "cueFile", "use cue sheet as markers")
	show.String(&cmd.Flags.CoverFile, "a", "artFile", "update album art")
	show.Bool(&cmd.Flags.MetaSwitch, "M", "meta", "Remove meta")
	show.Bool(&cmd.Flags.ChapSwitch, "C", "chaps", "Remove chapters")
	show.Bool(&cmd.Flags.CoverSwitch, "A", "art", "Remove album art")
	flaggy.AttachSubcommand(show, 1)

	// Extract command
	extract := flaggy.NewSubcommand("extract")
	extract.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	extract.Bool(&cmd.Flags.MetaSwitch, "m", "meta", "extract meta")
	extract.Bool(&cmd.Flags.ChapSwitch, "c", "chaps", "extract chapters")
	extract.Bool(&cmd.Flags.CoverSwitch, "a", "art", "extract album art")
	flaggy.AttachSubcommand(extract, 1)
	// Extract alias
	x := flaggy.NewSubcommand("x")
	x.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	x.Bool(&cmd.Flags.MetaSwitch, "m", "meta", "extract meta")
	x.Bool(&cmd.Flags.ChapSwitch, "c", "chaps", "extract chapters")
	x.Bool(&cmd.Flags.CoverSwitch, "a", "art", "extract album art")
	flaggy.AttachSubcommand(x, 1)

	// Update command
	update := flaggy.NewSubcommand("update")
	update.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	update.String(&cmd.Flags.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	update.String(&cmd.Flags.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(update, 1)
	// Update alias
	u := flaggy.NewSubcommand("u")
	u.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	u.String(&cmd.Flags.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	u.String(&cmd.Flags.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(u, 1)

	flaggy.Parse()

	// Input

	// Handle flags
	//if cmd.Overwrite {
	//  cmd.OverWrite()
	//}

	//if cmd.Verbose {
	//  cmd.Args().Verbose()
	//}

	switch {
	case join.Used, j.Used:
		//cmd.Join()
	case show.Used:
		cmd.Show()
		//cmd.Show()
	case cut.Used, c.Used:
		//cmd.Cut(cmd.Start, cmd.End, cmd.ChapNo).Run()
	case split.Used, s.Used:
		//err := cmd.Split()
		//if err != nil {
			//log.Fatal(err)
		//}
	case rm.Used, remove.Used:
		//cmd.Remove()
	case extract.Used, x.Used:
		cmd.Extract()
	case update.Used, u.Used:
		//cmd.Update()
	}
}
