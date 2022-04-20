package main

import (
	"fmt"
	//"strings"
	"bytes"
	"log"

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

	split := splitCmd()
	flaggy.AttachSubcommand(split, 1)

	convert := convertCmd()
	flaggy.AttachSubcommand(convert, 1)

	cue := cueCmd()
	flaggy.AttachSubcommand(cue, 1)

	rm := rmCmd()
	flaggy.AttachSubcommand(rm, 1)
	rmChaps := rmChapsCmd()
	rm.AttachSubcommand(rmChaps, 1)
	rmCover := rmCoverCmd()
	rm.AttachSubcommand(rmCover, 1)
	rmMeta := rmMetaCmd()
	rm.AttachSubcommand(rmMeta, 1)

	embed := embedCmd()
	flaggy.AttachSubcommand(embed, 1)
	embedChaps := embedChapsCmd()
	embed.AttachSubcommand(embedChaps, 1)
	embedCover := embedCoverCmd()
	embed.AttachSubcommand(embedCover, 1)
	embedMeta := embedMetaCmd()
	embed.AttachSubcommand(embedMeta, 1)

	extract := extractCmd()
	flaggy.AttachSubcommand(extract, 1)
	extractChaps := extractChapsCmd()
	extract.AttachSubcommand(extractChaps, 1)
	extractCover := extractCoverCmd()
	extract.AttachSubcommand(extractCover, 1)
	extractMeta := extractMetaCmd()
	extract.AttachSubcommand(extractMeta, 1)

	flaggy.Parse()

	// Setup command
	cmd := fftools.NewCmd().Profile(profile)

	// Input
	if posInput != "" {
		cmd.In(posInput)
	}
	for _, in := range input {
		cmd.In(in)
	}

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

	if cue.Used {
		fmt.Println("split")
	}

	if join.Used {
		c := cmd.Join(ext).Cmd()
		var out bytes.Buffer
		c.Stderr = &out
		err := c.Run()
		fmt.Printf("%q\n", out.String())
		log.Printf("Command finished with error: %v", err)
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

	fmt.Println(cmd.Cmd().String())
}

func joinCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("join")
	cmd.AddPositionalValue(&ext, "extension", 1, true, "specify the extension of the files")
	return cmd
}

func splitCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("split")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func cueCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("cue")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func convertCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("convert")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func rmCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("rm")
	return cmd
}

func rmCoverCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("cover")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func rmMetaCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("meta")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func rmChapsCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("chaps")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func extractCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("extract")
	return cmd
}

func extractMetaCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("meta")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func extractCoverCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("cover")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func extractChapsCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("chaps")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func embedCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("embed")
	return cmd
}

func embedMetaCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("meta")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func embedCoverCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("cover")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

func embedChapsCmd() *flaggy.Subcommand {
	cmd := flaggy.NewSubcommand("chaps")
	cmd.AddPositionalValue(&posInput, "input", 1, true, "input for the command")
	return cmd
}

