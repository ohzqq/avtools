package avtools

import (
	"github.com/ohzqq/avtools/timestamp"
	"golang.org/x/exp/constraints"
)

type Metadata interface {
	Chapters() []ChapterMeta
	Tags() map[string]string
	Streams() []map[string]string
}

type Meta struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
	Chapters []*Chapter
	streams  []map[string]string
}

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
