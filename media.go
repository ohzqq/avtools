package avtools

import (
	"encoding/json"
	"fmt"
	"log"
)

type Media struct {
	input string
	Meta
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
	Streams  []*Stream
	Chapters []*Chapter
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
}

func NewMedia(input string) *Media {
	media := Media{
		input: input,
	}
	//media.Probe()
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

func (m *Media) Probe() *Media {
	info := Probe(m.input)

	var raw map[string]json.RawMessage
	err := json.Unmarshal(info, &raw)
	if err != nil {
		log.Fatal(err)
	}

	var meta Meta
	err = json.Unmarshal(raw["format"], &meta)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(raw["streams"], &meta.Streams)
	if err != nil {
		log.Fatal(err)
	}

	meta.Chapters = UnmarshalChapters(raw["chapters"])

	m.Meta = meta

	fmt.Printf("%+V\n", meta)

	return m
}

func UnmarshalChapters(data json.RawMessage) []*Chapter {
	var chaps []chapterEntry
	err := json.Unmarshal(data, &chaps)
	if err != nil {
		log.Fatal(err)
	}

	var chapters []*Chapter
	for _, ch := range chaps {
		chapters = append(chapters, ch.ToChapter())
	}

	return chapters
}
