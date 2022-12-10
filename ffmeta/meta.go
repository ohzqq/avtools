package ffmeta

import (
	"math"
	"strconv"

	"github.com/ohzqq/avtools/chap"
)

type Meta struct {
	chap.Chapters `json:"-"`
	name          string
	Streams       []*Stream `json:"streams"`
	Format        `json:"format"`
	Chaps         []Chapter `json:"chapters"`
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string   `json:"filename"`
	Dur      Duration `json:"duration"`
	duration chap.Time
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

type Duration string

func NewFFmeta() *Meta {
	return &Meta{
		Chapters: chap.NewChapters(),
		Format:   Format{Tags: make(map[string]string)},
	}
}

func (ff Meta) Duration() chap.Time {
	t := chap.ParseStr(string(ff.Dur))
	return t
}

func (ff *Meta) SetTag(key, val string) *Meta {
	ff.Format.Tags[key] = val
	return ff
}

func (ff Meta) GetTags() map[string]string {
	return ff.Format.Tags
}

func (ff Meta) HasAudio() bool {
	return len(ff.AudioStreams()) > 0
}

func (ff Meta) AudioStreams() []*Stream {
	var streams []*Stream
	for _, s := range ff.Streams {
		if s.CodecType == "audio" {
			streams = append(streams, s)
		}
	}
	return streams
}

func (ff Meta) HasVideo() bool {
	return len(ff.VideoStreams()) > 0
}

func (ff Meta) VideoStreams() []*Stream {
	var streams []*Stream
	for _, s := range ff.Streams {
		if s.CodecType == "video" {
			streams = append(streams, s)
		}
	}
	return streams
}

func (d Duration) String() string {
	return string(d)
}

func (d Duration) Int() int {
	return int(math.Round(d.Float()))
}

func (d Duration) Float() float64 {
	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		return 0
	}
	return f
}
