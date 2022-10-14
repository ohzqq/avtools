package avtools

import (
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

type Media struct {
	files map[string]*FileFormat
	//Input         *FileFormat
	CueChaps bool
	outMeta  *FileFormat
	json     []byte
}

func NewMedia(input string) *Media {
	media := Media{
		files: make(map[string]*FileFormat),
	}
	media.SetFile("input", input)
	return &media
}

func (m Media) Meta() *MediaMeta {
	meta := m.GetFile("input").Meta()

	if m.HasFFmeta() {
		ff := m.GetFile("ffmeta")
		if ff.Meta().HasChapters() {
			meta.SetChapters(ff.meta.Chapters)
		}
		meta.SetTags(ff.meta.Format.Tags)
	}

	if m.HasCue() {
		meta.SetChapters(m.GetFile("cue").meta.Chapters)
	}

	return meta
}

func (m *Media) SafeName() string {
	var filename string
	if f := m.Meta().Format.Filename; f != "" {
		filename = f
	}
	if f := m.GetFile("input"); f != nil {
		filename = f.Name()
	}
	filename = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	return slug.Make(filename)
}

func (m *Media) GetFile(file string) *FileFormat {
	return m.files[file]
}

func (m *Media) SetFile(name, f string) *Media {
	file := NewFormat(name)
	if f != "" {
		file.SetFile(f)
		if name != "cover" {
			file.Parse()
		}
	}
	m.files[name] = file
	return m
}

func (m Media) HasCue() bool {
	return m.files["cue"] != nil
}

func (m Media) HasCover() bool {
	return m.files["cover"] != nil
}

func (m Media) HasFFmeta() bool {
	return m.files["ffmeta"] != nil
}

func (m Media) HasChapters() bool {
	if len(m.Meta().Chapters) != 0 {
		return true
	}
	return false
}

func (m Media) HasVideo() bool {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "video" {
			return true
		}
	}
	return false
}

func (m Media) HasAudio() bool {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "audio" {
			return true
		}
	}
	return false
}

func (m *Media) VideoCodec() string {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "video" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) AudioCodec() string {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "audio" {
			return stream.CodecName
		}
	}
	return ""
}

//func (m *Media) HasStreams() bool {
//  if m.HasMeta() && m.Meta.Streams != nil {
//    return true
//  }
//  return false
//}

//func (m *Media) HasFormat() bool {
//  if m.HasMeta() && m.Meta.Format != nil {
//    return true
//  }
//  return false
//}

func (m *Media) Duration() string {
	return secsToHHMMSS(m.Meta().Format.Duration)
}
