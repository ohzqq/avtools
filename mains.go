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
	args := avtools.NewArgs()

	// Flags
	flaggy.StringSlice(&cmd.InputSlice, "i", "input", "input")
	flaggy.String(&args.Name, "o", "output", "Set output")
	flaggy.Bool(&args.Overwrite, "y", "yes", "overwrite")
	flaggy.Bool(&args.Verbose, "v", "verbose", "print command")
	flaggy.String(&args.Profile, "p", "profile", "designate profile")

	// Subcommands

	// Join command
	join := flaggy.NewSubcommand("join")
	join.AddPositionalValue(&args.Extension, "extension", 1, true, "specify the extension of the files")
	flaggy.AttachSubcommand(join, 1)
	// Join alias
	j := flaggy.NewSubcommand("j")
	j.AddPositionalValue(&args.Extension, "extension", 1, true, "specify the extension of the files")
	flaggy.AttachSubcommand(j, 1)

	// Cut command
	cut := flaggy.NewSubcommand("cut")
	cut.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	cut.String(&args.MetaFile, "m", "meta", "use ffmetadata file as markers")
	cut.String(&args.CueFile, "c", "cue", "use cue sheet as markers")
	cut.String(&args.Start, "ss", "start", "the start time")
	cut.String(&args.End, "to", "end", "the end time")
	cut.Int(&args.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(cut, 1)
	// Cut alias
	c := flaggy.NewSubcommand("c")
	c.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	c.String(&args.MetaFile, "m", "meta", "use ffmetadata file as markers")
	c.String(&args.CueFile, "c", "cue", "use cue sheet as markers")
	c.String(&args.Start, "ss", "start", "the start time")
	c.String(&args.End, "to", "end", "the end time")
	c.Int(&args.ChapNo, "n", "num", "chapter to cut")
	flaggy.AttachSubcommand(c, 1)

	// Split command
	split := flaggy.NewSubcommand("split")
	split.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	split.String(&args.CueFile, "c", "cue", "use cue sheet as markers")
	split.String(&args.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(split, 1)
	// Split alias
	s := flaggy.NewSubcommand("s")
	s.AddPositionalValue(&args.Input, "input", 1, true, "input for the command")
	s.String(&args.CueFile, "c", "cue", "use cue sheet as markers")
	s.String(&args.MetaFile, "m", "meta", "use ffmetadata file as markers")
	flaggy.AttachSubcommand(s, 1)

	// Convert command
	//convert := flaggy.NewSubcommand("convert")
	//convert.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	//flaggy.AttachSubcommand(convert, 1)

	// Remove command
	remove := flaggy.NewSubcommand("remove")
	remove.AddPositionalValue(&args.Input, "input", 1, true, "input for the command")
	remove.Bool(&args.MetaFlag, "m", "meta", "Remove meta")
	remove.Bool(&args.ChapFlag, "c", "chaps", "Remove chapters")
	remove.Bool(&args.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(remove, 1)
	// Remove alias
	rm := flaggy.NewSubcommand("rm")
	rm.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	rm.Bool(&args.MetaFlag, "m", "meta", "Remove meta")
	rm.Bool(&args.ChapFlag, "c", "chaps", "Remove chapters")
	rm.Bool(&args.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(rm, 1)

	// Show command
	show := flaggy.NewSubcommand("show")
	show.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	show.Bool(&args.MetaFlag, "m", "meta", "Remove meta")
	show.Bool(&args.ChapFlag, "c", "chaps", "Remove chapters")
	show.Bool(&args.CoverFlag, "a", "art", "Remove album art")
	flaggy.AttachSubcommand(show, 1)

	// Extract command
	extract := flaggy.NewSubcommand("extract")
	extract.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	extract.Bool(&args.MetaFlag, "m", "meta", "extract meta")
	extract.Bool(&args.ChapFlag, "c", "chaps", "extract chapters")
	extract.Bool(&args.CoverFlag, "a", "art", "extract album art")
	flaggy.AttachSubcommand(extract, 1)
	// Extract alias
	x := flaggy.NewSubcommand("x")
	x.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	x.Bool(&args.MetaFlag, "m", "meta", "extract meta")
	x.Bool(&args.ChapFlag, "c", "chaps", "extract chapters")
	x.Bool(&args.CoverFlag, "a", "art", "extract album art")
	flaggy.AttachSubcommand(x, 1)

	// Update command
	update := flaggy.NewSubcommand("update")
	update.AddPositionalValue(&cmd.Input, "input", 1, true, "input for the command")
	update.String(&args.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	update.String(&args.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(update, 1)
	// Update alias
	u := flaggy.NewSubcommand("u")
	u.AddPositionalValue(&args.Input, "input", 1, true, "input for the command")
	u.String(&args.MetaFile, "m", "meta", "update meta")
	//embed.Bool(&xChaps, "c", "chaps", "update chapters")
	u.String(&args.CoverFile, "a", "art", "update album art")
	flaggy.AttachSubcommand(u, 1)

	flaggy.Parse()

	// Input
	if cmd.Input != "" {
		//cmd.Media.Input(cmd.Input)
	}

	cmd.CliArgs = args

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
		//cmd.Extract()
	case update.Used, u.Used:
		//cmd.Update()
	}
}
