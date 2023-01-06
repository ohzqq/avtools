package avtools

type Media struct {
	input string
	Meta
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
	Chapters []*Chapter
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
}

func NewMedia(input string) *Media {
	media := Media{
		input: input,
	}
	return &media
}

func (m *Media) SetMeta(meta Metadata) *Media {
	if tags := meta.Tags(); tags != nil {
		m.Tags = tags
	}
	if chaps := meta.Chapters(); len(chaps) > 0 {
		for _, chap := range chaps {
			ch := NewChapter(
				chap.Start(),
				chap.End(),
				chap.Timebase(),
			)
			ch.title = chap.Title()
			m.Chapters = append(m.Chapters, ch)
		}
	}
	//if streams := meta.Streams(); len(streams) > 0 {
	m.streams = meta.Streams()
	//}
	return m
}
