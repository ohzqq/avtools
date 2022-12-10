package ffmeta

import (
	"math"
	"strconv"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/ffprobe"
)

type Meta struct {
	Chapters     chap.Chapters `json:"chapters"`
	ffprobe.Meta `json:"-"`
	name         string
}

type Duration string

func NewFFmeta() *Meta {
	return &Meta{
		Chapters: chap.NewChapters(),
	}
}

func (ff *Meta) AddChapter(ch *chap.Chapter) {
	ff.Chapters.Chapters = append(ff.Chapters.Chapters, ch)
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

func (ff Meta) AudioStreams() []*ffprobe.Stream {
	var streams []*ffprobe.Stream
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

func (ff Meta) VideoStreams() []*ffprobe.Stream {
	var streams []*ffprobe.Stream
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
