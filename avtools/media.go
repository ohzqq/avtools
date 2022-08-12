package avtools

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type Media struct {
	Input    *FileFormat
	Cue      *FileFormat
	FFmeta   *FileFormat
	Json     *FileFormat
	Cover    *FileFormat
	CueChaps bool
	json     []byte
}

func NewMedia(input string) *Media {
	media := Media{
		Input: NewFormat(input),
	}
	if media.Input.IsAudio() {
		media.Input.parse = EmbeddedJsonMeta
		media.Input.render = MarshalJson
		media.Input.Parse()
	}
	return &media
}

func (m Media) Meta() *MediaMeta {
	meta := m.Input.meta
	if m.HasFFmeta() {
		meta = m.FFmeta.meta
	}
	fmt.Printf("%+V\n", meta)
	if m.HasCue() {
		meta.SetChapters(m.Cue.meta.Chapters)
	}
	return meta
}

func (m *Media) AddCover(file string) *Media {
	m.Cover = NewFormat(file)
	return m
}

func (m *Media) AddFFmeta(file string) *Media {
	m.FFmeta = NewFormat(file)
	m.FFmeta.parse = LoadFFmetadataIni
	m.FFmeta.render = RenderTmpl
	m.FFmeta.tmpl = template.Must(template.New("ffmeta").Funcs(funcs).Parse(ffmetaTmpl))
	m.FFmeta.Parse()
	return m
}

func (m *Media) AddCue(file string) *Media {
	m.Cue = NewFormat(file)
	m.Cue.parse = LoadCueSheet
	m.Cue.render = RenderTmpl
	m.Cue.tmpl = template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl))
	m.Cue.Parse()
	return m
}

//func (m *Media) ConvertTo(kind string) *FileFormat {
//  f := m.GetFormat(kind)
//  f.render(f)

//  switch kind {
//  case "json", ".json":
//    f.data = MarshalJson(f)
//  case "ffmeta", "ini", ".ini":
//    f.data = RenderTmpl(f)
//  case "cue", ".cue":
//    if len(f.meta.Chapters) == 0 {
//      log.Fatal("No chapters")
//    }
//    f.data = RenderTmpl(f)
//        return fmt.Render("cue")
//  }
//  return f
//}

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

func (m Media) HasCue() bool {
	return m.Cue != nil
}

func (m Media) HasCover() bool {
	return m.Cover != nil
}

func (m Media) HasFFmeta() bool {
	return m.FFmeta != nil
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
