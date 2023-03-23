package probe

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Meta struct {
	StreamEntry  []map[string]any `json:"streams"`
	Format       Format           `json:"format"`
	ChapterEntry []Chapter        `json:"chapters"`
}

type Format struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

func Load(input string) (Meta, error) {
	args := ffmpeg.MergeKwArgs(probeArgs)
	info, err := ffmpeg.ProbeWithTimeoutExec(input, 0, args)
	if err != nil {
		return Meta{}, err
	}
	//fmt.Println(info)

	data := []byte(info)

	var meta Meta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return Meta{}, err
	}

	return meta, nil
}

func DumpFFMeta(file string) *ff.Cmd {
	cmd := ff.New()
	cmd.In(file, ffmpeg.KwArgs{"y": ""})
	dir, name := filepath.Split(file)
	ext := filepath.Ext(name)
	name = "ffmeta-" + strings.TrimSuffix(name, ext)
	name = filepath.Join(dir, name)
	cmd.Output.Pad("").Name(name).Ext(".ini")
	cmd.Output.Set("f", "ffmetadata")
	return cmd.Compile()
}

func (m Meta) Chapters() []avtools.ChapterMeta {
	var chaps []avtools.ChapterMeta
	for _, ch := range m.ChapterEntry {
		chaps = append(chaps, ch)
	}
	return chaps
}

func (m Meta) Streams() []map[string]string {
	var streams []map[string]string
	for _, stream := range m.StreamEntry {
		meta := make(map[string]string)
		for key, raw := range stream {
			switch val := raw.(type) {
			case float64:
				meta[key] = strconv.Itoa(int(val))
			case string:
				meta[key] = val
			case map[string]any:
				if key == "disposition" {
					if val["attached_pic"].(float64) == 0 {
						meta["cover"] = "false"
					} else {
						meta["cover"] = "true"
					}
				}
			}
		}
		streams = append(streams, meta)
	}
	return streams
}

func (m Meta) Tags() map[string]string {
	m.Format.Tags["filename"] = m.Format.Filename
	m.Format.Tags["duration"] = m.Format.Dur
	m.Format.Tags["size"] = m.Format.Size
	m.Format.Tags["bit_rate"] = m.Format.BitRate
	return m.Format.Tags
}

var probeArgs = []ffmpeg.KwArgs{
	ffmpeg.KwArgs{"show_chapters": ""},
	ffmpeg.KwArgs{"pretty": ""},
	//ffmpeg.KwArgs{"select_streams": "a"},
	ffmpeg.KwArgs{"show_entries": "stream:format=filename, start_time, duration, size, bit_rate:format_tags"},
	ffmpeg.KwArgs{"of": "json"},
}
