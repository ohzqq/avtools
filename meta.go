package avtools

import (
	"github.com/ohzqq/avtools/timestamp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/exp/constraints"
)

type Meta struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
	Streams  []*Stream
	Chapters []*Chapter
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	stream    *ffmpeg.Stream
}

type Chapters []*Chapter

type Chapter struct {
	start timestamp.Time
	end   timestamp.Time
	base  timestamp.Timebase
	Title string
	title string
}

type Num interface {
	constraints.Integer | constraints.Float
}

type ChapterMeta interface {
	Start() float64
	End() float64
	Timebase() float64
	Title() string
}

func NewerChapter(meta ChapterMeta) *Chapter {
	return &Chapter{
		base:  timestamp.Timebase(meta.Timebase()),
		start: timestamp.NewerTimeStamp(meta.Start(), meta.Timebase()),
		end:   timestamp.NewerTimeStamp(meta.End(), meta.Timebase()),
		title: meta.Title(),
	}
}
