package avtools

type Media struct {
	Filename string
	Dur      string
	Size     string
	BitRate  string
	Tags     map[string]string
	Chapters []*Chapter
	streams  []map[string]string
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
}

type Metadata interface {
	Chapters() []ChapterMeta
	Tags() map[string]string
	Streams() []map[string]string
}

func NewMedia(input string) *Media {
	media := Media{
		Filename: input,
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
