package media

import (
	"log"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	"github.com/ohzqq/avtools/meta"
)

type Media struct {
	*avtools.Media
	streams     []Stream
	Input       File
	Output      File
	Ini         File
	Cue         File
	Cover       File
	HasCover    bool
	MetaChanged bool
}

type Stream struct {
	CodecType string
	CodecName string
	Index     string
	IsCover   bool
}

func New(input string) *Media {
	m := avtools.NewMedia(input)
	med := &Media{
		Media: m,
		Input: NewFile(input),
	}
	med.Output = File{FileName: med.Input.NewName()}
	med.Probe()

	return med
}

func (m *Media) LoadMeta(name string) *Media {
	file := NewFile(name)
	switch {
	case file.IsFFMeta():
		ini := meta.LoadIni(file.Abs)
		m.Media.SetMeta(ini)
		m.Ini = file
	case file.IsCue():
		cue := meta.LoadCueSheet(file.Abs)
		m.Media.SetMeta(cue)
		dur := m.GetTag("duration")
		last := m.Chapters()[len(m.Chapters())-1]
		last.End = avtools.Timestamp(avtools.ParseStamp(dur))
		m.Cue = file
	}
	m.MetaChanged = true
	return m
}

func (m Media) Command() ff.Cmd {
	cmd := ff.New()
	cmd.In(m.Input.Abs)
	return cmd
}

func (m Media) HasChapters() bool {
	return len(m.Chapters()) > 0
}

func (m Media) AudioStreams() []Stream {
	var streams []Stream
	for _, stream := range m.streams {
		if stream.CodecType == "audio" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (m Media) VideoStreams() []Stream {
	var streams []Stream
	for _, stream := range m.streams {
		if stream.CodecType == "video" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func IsPlainText(mtype string) bool {
	if strings.Contains(mtype, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}
