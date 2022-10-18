package ffmeta

import (
	"math"
	"strconv"

	"github.com/ohzqq/avtools/chap"
)

type FFmeta struct {
	chap.Chapters
	name    string
	Streams []*Stream
	Format  `json:"format"`
	Chaps   []Chapter `json:"chapters"`
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	Dur      Duration `json:"duration"`
	Size     string
	BitRate  string `json:"bit_rate"`
	Tags     map[string]string
}

type Duration string

func NewFFmeta() *FFmeta {
	return &FFmeta{Chapters: chap.NewChapters()}
}

func (ff FFmeta) Duration() Duration {
	return ff.Dur
}

func (ff FFmeta) HasAudio() bool {
	return len(ff.AudioStreams()) > 0
}

func (ff FFmeta) AudioStreams() []*Stream {
	var streams []*Stream
	for _, s := range ff.Streams {
		if s.CodecType == "audio" {
			streams = append(streams, s)
		}
	}
	return streams
}

func (ff FFmeta) HasVideo() bool {
	return len(ff.VideoStreams()) > 0
}

func (ff FFmeta) VideoStreams() []*Stream {
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
