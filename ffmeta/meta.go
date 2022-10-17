package ffmeta

import "github.com/ohzqq/avtools/chap"

type FFmeta struct {
	chap.Chapters
	name string
	Tags map[string]string
}

func NewFFmeta() *FFmeta {
	return &FFmeta{Chapters: chap.NewChapters()}
}

func LoadJson(d []byte) ffmeta {
}
