package yygif

import (
	"github.com/ohzqq/avtools/media"
	"github.com/ohzqq/avtools/meta"
)

type Gif struct {
	*media.Media
}

func MkGifs(input string) Gif {
	meta := meta.LoadIni(input)
	vid := meta.Tags()["title"]
	src := media.New(vid)
	src.SetMeta(meta)
	return Gif{Media: src}
}
