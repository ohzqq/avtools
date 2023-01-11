package yygif

import (
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	"github.com/ohzqq/avtools/media"
	"github.com/ohzqq/avtools/meta"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Gif struct {
	*avtools.Chapter
	Input *media.Media
	Crop  string
}

func LoadGifMeta(input string) *media.Media {
	meta := meta.LoadIni(input)
	src := avtools.NewMedia().SetMeta(meta)
	vid := meta.Tags()["title"]
	return &media.Media{
		Media:   src,
		Input:   media.NewFile(vid),
		Profile: "gif",
	}
}

func MkGif(input string, ch *avtools.Chapter) *ff.Cmd {
	in := media.NewFile(input)
	cmd := ff.New("gif")
	inKwargs := ffmpeg.KwArgs{
		"ss": ch.Start.String(),
		"to": ch.End.String(),
	}
	cmd.In(in.Abs, inKwargs)

	filters := ff.GetProfile("gif").Filters
	if c, ok := ch.Tags["crop"]; ok {
		crop := strings.Split(c, ":")
		filters.Set("crop", crop...)
	}
	cmd.Filters = filters

	out := in.NewName()
	if ch.Title != "" {
		cmd.Output.Pad("")
		out.Name = ch.Title
	} else {
		out.Suffix("-")
	}

	cmd.Output.Name(out.Join())

	return &cmd
}
