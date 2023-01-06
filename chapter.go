package avtools

import (
	"fmt"

	"github.com/ohzqq/avtools/timestamp"
)

type Number interface {
	int | int32 | int64 | float32 | float64
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
		base:  timestamp.Timebase(base),
		start: timestamp.NewerTimeStamp(start, base),
		end:   timestamp.NewerTimeStamp(end, base),
	}
}

func (ch Chapter) Start() timestamp.Time {
	return ch.start
}

func (ch Chapter) End() timestamp.Time {
	return ch.end
}

func (ch Chapter) Timebase() timestamp.Timebase {
	return ch.base
}

func (ch Chapter) Dur() (timestamp.Time, error) {
	if ch.end.Duration == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.Duration - ch.start.Duration
	stamp := timestamp.NewerTimeStamp(t, float64(ch.base))
	return stamp, nil
}
