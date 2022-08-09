package avtools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Media struct {
	Meta     *MediaMeta
	Format   *FileFormat
	Formats  []*FileFormat
	File     string
	Path     string
	Dir      string
	Ext      string
	mimetype string
	CueChaps bool
	json     []byte
}

func NewMedia(input string) *Media {
	mime.AddExtensionType(".ini", "text/plain")
	mime.AddExtensionType(".cue", "text/plain")
	m := Media{}

	abs, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}

	m.Path = abs
	m.File = filepath.Base(input)
	m.Dir = filepath.Dir(input)
	m.Ext = filepath.Ext(input)
	m.mimetype = mime.TypeByExtension(m.Ext)
	m.Formats = []*FileFormat{
		&FileFormat{
			name:   "cue",
			ext:    ".cue",
			parse:  LoadCueSheet,
			render: RenderTmpl,
			tmpl:   template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
		},
		&FileFormat{
			name:   "ffmeta",
			ext:    ".ini",
			parse:  LoadFFmetadataIni,
			render: RenderTmpl,
			tmpl:   template.Must(template.New("ffmeta").Funcs(funcs).Parse(ffmetaTmpl)),
		},
		&FileFormat{
			name:   "json",
			ext:    ".json",
			parse:  JsonMeta,
			render: MarshalJson,
		},
	}

	return &m
}

func (m *Media) GetFormat(f string) *FileFormat {
	for _, format := range m.Formats {
		if f == format.ext {
			return format
		}
	}
	return m.Formats[2]
}

func (m *Media) AddFileFormat(f string) *Media {
	format := m.GetFormat(path.Ext(f))
	format.file = f
	format.Parse()
	return m
}

func (m *Media) JsonMeta() *Media {
	cmd := NewFFprobeCmd(m.Path)
	cmd.entries = ffProbeMeta
	cmd.showChaps = true
	cmd.format = "json"
	m.json = cmd.Parse().Run()
	return m
}

func (m *Media) Unmarshal() *Media {
	err := json.Unmarshal(m.json, &m.Meta)
	if err != nil {
		fmt.Println("help")
	}
	m.Meta.Tags = m.Meta.Format.Tags
	return m
}

func (m *Media) Print() {
	fmt.Println(string(m.json))
}

//func (m *Media) RenderFFChaps() string {
//  if m.Meta == nil {
//    m.JsonMeta().Unmarshal()
//  }

//  lastCh := m.Meta.Chapters[len(m.Meta.Chapters)-1]
//  lastCh.End = m.Meta.Format.DurationSecs(lastCh.TimebaseFloat())

//  fmt, err := GetFormat("cue")
//  if err != nil {
//    log.Println(err)
//  }
//  //fmt.SetMeta(m.Meta).ConvertTo("ffmeta")

//  return fmt.String()
//}

func (m *Media) SetMeta(meta *MediaMeta) *Media {
	m.Meta.Format.Tags = meta.Format.Tags
	m.Meta.Tags = meta.Tags
	m.Meta.Chapters = meta.Chapters
	return m
}

func (m Media) IsPlainText() bool {
	if strings.Contains(m.mimetype, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

func (m Media) IsMeta() bool {
	m.IsPlainText()

	contents, err := os.Open(m.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	scanner := bufio.NewScanner(contents)
	line := 0
	for scanner.Scan() {
		if line == 0 && scanner.Text() == ";FFMETADATA1" {
			return true
		} else {
			log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
		}
	}
	return false
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

func (m Media) IsImage() bool {
	if strings.Contains(m.mimetype, "image") {
		return true
	} else {
		log.Fatalln("this switch requires an image file")
	}
	return false
}

func (m *Media) SetChapters(ch Chapters) {
	m.Meta.Chapters = ch
}

func (m *Media) HasChapters() bool {
	if m.HasMeta() && len(m.Meta.Chapters) != 0 {
		return true
	}
	return false
}

func (m *Media) HasVideo() bool {
	if m.HasMeta() {
		for _, stream := range m.Meta.Streams {
			if stream.CodecType == "video" {
				return true
			}
		}
	}
	return false
}

func (m *Media) HasAudio() bool {
	if m.HasMeta() {
		for _, stream := range m.Meta.Streams {
			if stream.CodecType == "audio" {
				return true
			}
		}
	}
	return false
}

func (m *Media) VideoCodec() string {
	if m.HasMeta() {
		for _, stream := range m.Meta.Streams {
			if stream.CodecType == "video" {
				return stream.CodecName
			}
		}
	}
	return ""
}

func (m *Media) AudioCodec() string {
	if m.HasMeta() {
		for _, stream := range m.Meta.Streams {
			if stream.CodecType == "audio" {
				return stream.CodecName
			}
		}
	}
	return ""
}

func (m *Media) HasMeta() bool {
	if m.Meta != nil {
		return true
	}
	return false
}

func (m *Media) HasStreams() bool {
	if m.HasMeta() && m.Meta.Streams != nil {
		return true
	}
	return false
}

func (m *Media) HasFormat() bool {
	if m.HasMeta() && m.Meta.Format != nil {
		return true
	}
	return false
}

func (m *Media) Duration() string {
	return secsToHHMMSS(m.Meta.Format.Duration)
}
