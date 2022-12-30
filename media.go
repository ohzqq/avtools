package avtools

type Media struct {
	input string
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
	*Meta
}

func NewMedia(input string) *Media {
	media := Media{
		input: input,
		Meta:  ReadEmbeddedMeta(input),
	}
	return &media
}
