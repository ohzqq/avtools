package avtools

import (
	"path/filepath"
	"log"
	"fmt"
	//"os"
	//"bytes"
	"encoding/json"
	//"strconv"
	//"strings"
)
var _ = fmt.Printf

type Media struct {
	Meta *MediaMeta
	File string
	Path string
	Dir string
	Ext string
	CueChaps bool
	json []byte
}

func NewMedia(input string) *Media {
	m := Media{}

	abs, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}

	m.Path = abs
	m.File = filepath.Base(input)
	m.Dir = filepath.Dir(input)
	m.Ext = filepath.Ext(input)

	return &m
}

func(m *Media) JsonMeta() []byte {
	cmd := NewFFprobeCmd(m.Path)
	cmd.entries = ffProbeMeta
	cmd.showChaps = true
	cmd.format = "json"
	e := cmd.Parse()

	//cmd.Media.ffChaps = true
	return e.Run()
}
func(cmd *Cmd) GetJsonMeta() []byte {
	cmd.ffprobe = true
	cmd.args = NewArgs()
	cmd.args.entries = ffProbeMeta
	cmd.args.showChaps = true
	cmd.args.format = "json"
	c := cmd.parseFFprobeArgs()

	//cmd.Media.ffChaps = true
	return c.Run()
}

func(cmd *Cmd) ParseJsonMeta() *Cmd {
	err := json.Unmarshal(cmd.GetJsonMeta(), &cmd.Media.Meta)
	if err != nil {
		fmt.Println("help")
	}
	fmt.Printf("%+V\n", string(cmd.GetJsonMeta()))


	return cmd
}

func(cmd *Cmd) PrintJsonMeta() {
	fmt.Println(string(cmd.GetJsonMeta()))
}

func(m *Media) SetMeta(meta *MediaMeta) *Media {
	m.Meta = meta
	return m
}

func (m *Media) SetChapters(ch []*Chapter) {
	m.Meta.Chapters = ch
}

func(m *Media) HasChapters() bool {
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
	if m.HasMeta() && m.Meta.Streams != nil{
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

