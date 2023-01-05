package avtools

import (
	"encoding/json"
	"fmt"
	"log"
)

type Media struct {
	input string
	FFmeta
	//Input    MediaFile
	//Files    RelatedFiles
	//cueSheet *cue.Sheet
}

func NewMedia(input string) *Media {
	media := Media{
		input: input,
	}
	media.Probe()
	return &media
}

func (m *Media) Probe() *Media {
	info := Probe(m.input)

	var raw map[string]json.RawMessage
	err := json.Unmarshal(info, &raw)
	if err != nil {
		log.Fatal(err)
	}

	var meta FFmeta
	err = json.Unmarshal(raw["format"], &meta)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(raw["streams"], &meta.Streams)
	if err != nil {
		log.Fatal(err)
	}

	meta.Chapters = UnmarshalChapters(raw["chapters"])

	m.FFmeta = meta

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
