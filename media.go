package avtools

type Media struct {
	input string
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
	*ProbeMeta
}

func NewMedia(input string) *Media {
	media := Media{
		input:     input,
		ProbeMeta: ReadEmbeddedMeta(input),
	}
	return &media
}
