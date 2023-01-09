package media

import (
	"bufio"
	"html/template"
	"log"
	"os"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	"github.com/ohzqq/avtools/meta"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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
			if line == 0 && scanner.Text() == ";FFMETADATA1" {
				ini := meta.LoadIni(file.Abs)
				m.Media.SetMeta(ini)
				m.Ini = file
				break
			} else {
				log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
			}
		}
	}
	return m
}

func (m Media) DumpIni() []byte {
	return meta.DumpIni(m)
}

func (m *Media) LoadCue(name string) *Media {
	file := NewFile(name)
	if IsPlainText(file.Mimetype) {
		cue := meta.LoadCueSheet(file.Abs)
		m.Media.SetMeta(cue)
		dur := m.GetTag("duration")
		last := m.Chapters()[len(m.Chapters())-1]
		last.End = avtools.Timestamp(avtools.ParseStamp(dur))
	}
	return m
}

func (m Media) DumpCue() []byte {
	return meta.DumpCueSheet(m.Input.Abs, m)
}

func (m *Media) Probe() *Media {
	p := meta.FFProbe(m.Input.Abs)
	m.Media.SetMeta(p)

	if len(m.Media.Streams()) > 0 {
		for _, stream := range m.Media.Streams() {
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
						m.HasCover = true
					}
				}
			}
			m.streams = append(m.streams, s)
		}
	}

	return m
}

func (m Media) DumpFFMeta() ff.Cmd {
	cmd := ff.New()
	cmd.In(m.Input.Abs, ffmpeg.KwArgs{"y": ""})
	name := m.Input.NewName()
	n := name.Prefix("ffmeta-").Join()
	cmd.Output.Pad("").Name(n).Ext(".ini")
	cmd.Output.Set("f", "ffmetadata")
	return cmd
}

var tmplFuncs = template.FuncMap{
	"inc": Inc,
}

func Inc(n int) int {
	return n + 1
}

const cueTmpl = `FILE "{{.Input.Name}}" {{.Input.Ext -}}
{{range $idx, $ch := .Chapters}}
TRACK {{inc $idx}} AUDIO
{{- if eq $ch.Title ""}}
  TITLE "Chapter {{inc $idx}}"
{{- else}}
  TITLE "{{$ch.Title}}"
{{- end}}
  INDEX 01 {{$ch.Start.MMSS}}
{{- end -}}`
