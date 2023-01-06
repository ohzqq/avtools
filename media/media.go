package media

import (
	"bufio"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/meta"
)

type Media struct {
	*avtools.Media
	Streams []Stream
	Input   File
	Output  File
	FFmeta  File
	Cue     File
	Cover   File
}

type Stream struct {
	CodecType string
	CodecName string
	Index     string
	IsCover   bool
}

func New(input string) *Media {
	m := avtools.NewMedia(input)
	return &Media{
		Media: m,
		Input: NewFile(input),
	}
}

func (m Media) AudioStreams() []Stream {
	var streams []Stream
	for _, stream := range m.Streams {
		if stream.CodecType == "audio" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (m Media) VideoStreams() []Stream {
	var streams []Stream
	for _, stream := range m.Streams {
		if stream.CodecType == "video" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (m *Media) LoadIni(name string) *Media {
	file := NewFile(name)
	if IsPlainText(file.Mimetype) {
		contents, err := os.Open(file.Abs)
		if err != nil {
			log.Fatal(err)
		}
		defer contents.Close()

		scanner := bufio.NewScanner(contents)
		line := 0
		for scanner.Scan() {
			if line == 0 && scanner.Text() == meta.FFmetaComment {
				ini := meta.LoadIni(file.Abs)
				m.SetMeta(ini)
				m.FFmeta = file
				break
			} else {
				log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
			}
		}
	}
	return m
}

func (m *Media) LoadCue(name string) *Media {
	file := NewFile(name)
	if IsPlainText(file.Mimetype) {
		cue := meta.LoadCueSheet(file.Abs)
		m.SetMeta(cue)
	}
	return m
}

func (m *Media) Probe() *Media {
	p := meta.FFProbe(m.Input.Abs)
	m.SetMeta(p)

	if len(m.Media.Streams) > 0 {
		for _, stream := range m.Media.Streams {
			s := Stream{}
			for key, val := range stream {
				switch key {
				case "codec_type":
					s.CodecType = val
				case "codec_name":
					s.CodecName = val
				case "index":
					s.Index = val
				case "cover":
					if val == "true" {
						s.IsCover = true
					}
				}
			}
			m.Streams = append(m.Streams, s)
		}
	}

	return m
}

func IsPlainText(mtype string) bool {
	if strings.Contains(mtype, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

type File struct {
	Abs      string
	Path     string
	Base     string
	Ext      string
	Name     string
	File     string
	Padding  string
	Mimetype string
	name     string
}

func NewFile(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	f := File{
		Base:    filepath.Base(abs),
		Ext:     filepath.Ext(abs),
		Abs:     abs,
		Padding: "%03d",
	}
	f.Mimetype = mime.TypeByExtension(f.Ext)
	f.Name = strings.TrimSuffix(abs, f.Ext)
	f.name = strings.TrimSuffix(f.Base, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func (f File) WithExt(e string) string {
	return filepath.Join(f.Path, f.name+e)
}

func (f File) AddSuffix(s string) string {
	name := f.name + s + f.Ext
	return filepath.Join(f.Path, name)
}

func (f File) Pad(i int) string {
	p := fmt.Sprintf(f.Padding, i)
	return f.AddSuffix(p)
}
