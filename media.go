package avtools

import "github.com/ohzqq/avtools/ffmeta"

type Media struct {
	input string
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
	*ffmeta.Meta
}

type Meta struct {
}
