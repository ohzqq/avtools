package yygif

import (
	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/media"
	"github.com/ohzqq/avtools/meta"
)

type Gif struct {
	*avtools.Chapter
	Crop string
}

func MkGifs(input string) *media.Media {
	meta := meta.LoadIni(input)
	src := avtools.NewMedia().SetMeta(meta)
	return src
}

func MkGif(input string, ch *avtools.Chapter) Gif {
	g := Gif{Chapter: ch}
	if crop, ok := ch.Tags["crop"]; ok {
		g.Crop = crop
	}
	return g
}
