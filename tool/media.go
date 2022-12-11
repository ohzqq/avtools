package tool

import (
	"encoding/json"
	"log"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/cue"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/ffprobe"
	"github.com/ohzqq/avtools/file"
	"golang.org/x/exp/slices"
)

type Media struct {
	Input    MediaFile
	Files    RelatedFiles
	cueSheet *cue.Sheet
	*ffmeta.Meta
}

func NewMedia(i string) *Media {
	media := Media{
		Input: MediaFile{File: file.New(i)},
		Files: make(RelatedFiles),
	}
	media.Meta = media.ReadEmbeddedMeta()
	return &media
}

func (m *Media) AddFile(name, path string) *Media {
	m.Files[name] = MediaFile{File: file.New(path)}
	return m
}

func (m Media) HasFFmeta() bool {
	return m.Files.Has("ffmeta")
}

func (m *Media) SetFFmeta(ff string) *Media {
	m.Files["ffmeta"] = MediaFile{File: file.New(ff)}
	return m
}

func (m Media) EachChapter() []*chap.Chapter {
	var ch []*chap.Chapter
	if len(m.Meta.Chapters.Chapters) > 0 {
		ch = m.Meta.Chapters.Each()
	}
	return ch
}

func (m Media) HasCue() bool {
	return m.Files.Has("cue")
}

func (m *Media) SetCue(c string) *Media {
	m.Files["cue"] = MediaFile{File: file.New(c)}
	return m
}

func (m Media) HasEmbeddedCover() bool {
	for _, v := range m.Meta.VideoStreams() {
		if v.CodecName == "mjpeg" || v.CodecName == "png" {
			return true
		}
	}
	return false
}

func (m Media) EmbeddedCoverExt() string {
	for _, v := range m.Meta.VideoStreams() {
		if v.CodecName == "mjpeg" {
			return ".jpg"
		}
		if v.CodecName == "png" {
			return ".png"
		}
	}
	return ""
}

func (m *Media) SetMeta() *Media {
	if m.HasFFmeta() {
		meta := m.ReadFFmeta()
		m.Meta.Format.Tags = meta.Format.Tags
		m.Meta.SetChapters(meta.Chapters)
	}

	if m.HasCue() {
		m.Meta.SetChapters(m.ReadCueSheet())
	}

	return m
}

func (m *Media) ReadEmbeddedMeta() *ffmeta.Meta {
	probe := ffprobe.New()
	probe.Input(m.Input.Abs).
		FormatEntry("filename", "start_time", "duration", "size", "bit_rate").
		StreamEntry("codec_type", "codec_name").
		Entry("format_tags").
		ShowChapters().
		Json()

	data, err := probe.Run()
	if err != nil {
		log.Fatal(err)
	}

	meta := ffmeta.NewFFmeta()
	meta.Meta = ffprobe.UnmarshalJSON(data)

	for _, c := range meta.Chaps {
		ch := chap.NewChapter().SetMeta(c)
		meta.AddChapter(ch)
	}

	return meta
}

func (m *Media) ReadCueSheet() chap.Chapters {
	var ch chap.Chapters
	if m.HasCue() {
		ch = chap.NewChapters().FromCue(m.Files.Get("cue").Abs)
	}
	return ch
}

func (m *Media) ReadFFmeta() *ffmeta.Meta {
	var ff *ffmeta.Meta
	if m.HasFFmeta() {
		ff = ffmeta.LoadIni(m.Files.Get("ffmeta").Abs)
	}
	return ff
}

func (m Media) AudioCodec() string {
	if m.HasAudio() {
		a := m.AudioStreams()
		return a[0].CodecType
	}
	return ""
}

var metaTags = []string{
	"album",
	"album_artist",
	"artist",
	"comment",
	"composer",
	"copyright",
	"date",
	"genre",
	"language",
	"lyrics",
	"title",
	"track",
}

func (m Media) DumpJson() []byte {
	meta := make(map[string]interface{})
	meta["mimetype"] = m.Input.Mimetype
	meta["file"] = m.Input.Abs
	meta["size"] = m.Meta.Size
	meta["duration"] = m.Duration().HHMMSS()
	meta["chapters"] = m.Chapters
	for key, val := range m.Meta.Tags {
		if slices.Contains(metaTags, key) {
			meta[key] = val
		}
	}
	data, err := json.Marshal(meta)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
