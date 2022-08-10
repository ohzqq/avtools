package avtools

import (
	"fmt"
	"log"
	"os"
)

type Media struct {
	*FileFormats
	*MediaMeta
	CueChaps bool
	json     []byte
}

func NewMedia(input string) *Media {
	m := Media{FileFormats: &FileFormats{}}
	m.AddFormat(input)
	//m.Ext = m.GetFormat("audio").Ext
	m.MediaMeta = m.GetFormat("audio").meta
	return &m
}

func (m Media) Meta() *MediaMeta {
	if m.HasFFmeta() {
		m.MediaMeta = m.GetFormat(".ini").meta
	}
	if m.HasCue() {
		m.MediaMeta.SetChapters(m.GetFormat(".cue").meta.Chapters)
	}
	return m.MediaMeta
}

func (m *Media) ConvertTo(kind string) *FileFormat {
	f := m.GetFormat(kind)
	f.render(f)

	//switch kind {
	//case "json", ".json":
	//  f.data = MarshalJson(f)
	//case "ffmeta", "ini", ".ini":
	//  f.data = RenderTmpl(f)
	//case "cue", ".cue":
	//  if len(f.meta.Chapters) == 0 {
	//    log.Fatal("No chapters")
	//  }
	//  f.data = RenderTmpl(f)
	//  //    return fmt.Render("cue")
	//}
	return f
}

func (m *Media) Print() {
	fmt.Println(string(m.json))
}

func (m *Media) FFmetaChapsToCue() {
	if !m.HasChapters() {
		log.Fatal("No chapters")
	}

	f, err := os.Create("chapters.cue")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := GetTmpl("cue")
	if err != nil {
		log.Println(err)
	}

	err = tmpl.Execute(f, m.Meta)
	if err != nil {
		log.Println("executing template:", err)
	}
}

//func (m *Media) SetChapters(ch Chapters) {
//  m.Meta.Chapters = ch
//}

func (m *Media) HasChapters() bool {
	if len(m.Meta().Chapters) != 0 {
		return true
	}
	return false
}

func (m *Media) HasVideo() bool {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "video" {
			return true
		}
	}
	return false
}

func (m *Media) HasAudio() bool {
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
