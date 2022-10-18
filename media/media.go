package media

import (
	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/cue"
	"github.com/ohzqq/avtools/ffmeta"
)

type Media struct {
	input    string
	files    RelatedFiles
	meta     *ffmeta.FFmeta
	cue      *cue.CueSheet
	chapters chap.Chapters
}

type RelatedFiles map[string]string

func New(i string) *Media {
	return &Media{
		input: i,
		files: make(RelatedFiles),
	}
}
