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

func NewChapter[N Number](times ...N) *Chapter {
	var base float64 = 1
	var start float64 = 0
	var end float64 = 0

	switch t := len(times); t {
	case 3:
		base = float64(times[2])
		fallthrough
	case 2:
		end = float64(times[1])
		fallthrough
	case 1:
		start = float64(times[0])
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
