package avtools

import (
	"path/filepath"
	"log"
	"fmt"
	"os"
	"bytes"
	"encoding/json"
	//"strconv"
	//"strings"
)
var _ = fmt.Printf

type Media struct {
	Overwrite bool
	File string
	Path string
	Dir string
	Ext string
	Meta *MediaMeta
}

func NewMedia() *Media {
	return &Media{}
}

func(m *Media) Input(input string) *Media {
	abs, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}

	m.Path = abs
	m.File = filepath.Base(input)
	m.Dir = filepath.Dir(input)
	m.Ext = filepath.Ext(input)
	m.ParseJsonMeta()

	return m
}

func(m *Media) PrintJsonMeta() {
	fmt.Println(string(m.GetJsonMeta()))
}

func(m *Media) ParseJsonMeta() {
	err := json.Unmarshal(m.GetJsonMeta(), &m.Meta)
	if err != nil {
		fmt.Println("help")
	}
}

func(m *Media) GetJsonMeta() []byte {
	ff := NewFFprobeCmd()
	ff.In(m.Path)
	ff.Args().Entries(ffProbeMeta).Chapters().Verbosity("error").Format("json")

	return ff.Run()
}

func(m *Media) WriteMeta() {
	WriteFFmetadata(m.Path)
}

func(m *Media) SetMeta(meta *MediaMeta) *Media {
	m.Meta = meta
	return m
}

func(m *Media) HasChapters() bool {
	if m.Meta != nil {
		if m.Meta.Chapters != nil {
			return true
		}
	}
	return false
}

func(m *Media) RenderFFChaps() {
	var f bytes.Buffer

	err := metaTmpl.ffchaps.ExecuteTemplate(&f, "ffchaps", m)
	if err != nil {
		log.Println("executing template:", err)
	}
	fmt.Println(f.String())
}

func(m *Media) FFmetaChapsToCue() {
	f, err := os.Create("chapters.cue")
	if err != nil {
		log.Fatal(err)
	}

	err = metaTmpl.cue.ExecuteTemplate(f, "cue", m)
	if err != nil {
		log.Println("executing template:", err)
	}
}

func (m *Media) SetChapters(ch []*Chapter) {
	m.Meta.Chapters = ch
}

func (m *Media) HasVideo() bool {
	for _, stream := range m.Meta.Streams {
		if stream.CodecType == "video" {
			return true
		}
	}
	return false
}

func (m *Media) HasAudio() bool {
	for _, stream := range m.Meta.Streams {
		if stream.CodecType == "audio" {
			return true
		}
	}
	return false
}

func (m *Media) VideoCodec() string {
	for _, stream := range m.Meta.Streams {
		if stream.CodecType == "video" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) AudioCodec() string {
	for _, stream := range m.Meta.Streams {
		if stream.CodecType == "audio" {
			return stream.CodecName
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
	if m.HasMeta() {
		if m.Meta.Streams != nil {
			return true
		}
	}
	return false
}

func (m *Media) hasFormat() bool {
	if m.Meta != nil {
		if m.Meta.Format != nil {
			return true
		}
	}
	return false
}

func (m *Media) Duration() string {
	return secsToHHMMSS(m.Meta.Format.Duration)
}

