package avtools

import "github.com/ohzqq/avtools/timestamp"

type Number interface {
	int | int32 | int64 | float32 | float64
}

type Chapter struct {
	start timestamp.Time
	end   timestamp.Time
	base  timestamp.Timebase
	Title string
}

func NewChapter[N Number](start, end N, base N) *Chapter {
	if base == 0 {
		base = 1
	}
	return &Chapter{
		base:  timestamp.Timebase(1),
		start: timestamp.NewerTimeStamp(start, base),
		end:   timestamp.NewerTimeStamp(end, base),
	}
}

func ChapterStart[N Number](ch *Chapter, num N) *Chapter {
	ch.start = timestamp.NewTimeStamp(num)
	return ch
}
