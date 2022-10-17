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
	Dur      duration `json:"duration"`
	Size     string
	BitRate  string `json:"bit_rate"`
	Tags     map[string]string
}

type duration string

func NewFFmeta() *FFmeta {
	return &FFmeta{Chapters: chap.NewChapters()}
}

func (ff *FFmeta) SetChapters(c chap.Chapters) *FFmeta {
	ff.Chapters = c
	return ff
}

func (ff FFmeta) LastChapterEnd() *chap.Chapter {
	ch := ff.LastChapter()
	if ch.End().Secs() == 0 && ff.Duration().Int() != 0 {
		ch.SetEnd(chap.NewChapterTime(ff.Duration().Float() * 1000))
	}
	return ch
}

func (ff FFmeta) Duration() duration {
	return ff.Dur
}

func (d duration) String() string {
	return string(d)
}

func (d duration) Int() int {
	return int(math.Round(d.Float()))
}

func (d duration) Float() float64 {
	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		return 0
	}
	return f
}
