package timestamp

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

type Time struct {
	Time float64
	Base float64
}

type Timebase float64

func NewChapterTime[N Number](num N) Time {
	return Time{
		Time: float64(num),
		Base: 1,
	}
}

func ParseStr(t string) Time {
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		log.Fatal(err)
	}
	return NewChapterTime(i)
}

func ParseNum[N Number](num N, dig int) string {
	return strconv.FormatFloat(float64(num), 'f', dig, 64)
}

func (ch *Time) SetTimebase(base float64) {
	ch.Base = base
}

func (ch Time) Secs() int {
	secs := ch.Time / ch.Base
	return int(math.Round(secs))
}

func (ch Time) Float() float64 {
	return ch.Time
}

func (ch Time) String() string {
	return ParseNum(ch.Time, 0)
}

func (ch Time) SecsString() string {
	return ParseNum(ch.Time/ch.Base, 3)
}

func (c Time) HHMMSS() string {
	secs := c.Secs()
	hh := secs / 3600
	mm := secs % 3600 / 60
	ss := secs % 3600 % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

func (c Time) MMSS() string {
	secs := c.Secs()
	mm := secs / 60
	ss := secs % 60
	return fmt.Sprintf("%02d:%02d", mm, ss)
}

func (t Timebase) String() string {
	return "1/" + ParseNum(t.Float(), 0)
}

func (t Timebase) Float() float64 {
	return float64(t)
}
