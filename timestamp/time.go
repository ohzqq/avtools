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
	Duration float64
	base     float64
	Base     Timebase
}

type Timebase float64

func (t Timebase) String() string {
	return "1/" + ParseNumber(t.Float(), 0)
}

func (t Timebase) Float() float64 {
	return float64(t)
}

func NewTimeStamp[N Number](num N) Time {
	return Time{
		Duration: float64(num),
		Base:     1,
	}
}

func NewerTimeStamp[N Number](t, base N) Time {
	return Time{
		Duration: float64(t),
		base:     float64(base),
	}
}

func ParseString(t string) Time {
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		log.Fatal(err)
	}
	return NewTimeStamp(i)
}

func StringToFloat(t string) float64 {
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func ParseNumber[N Number](num N, dig int) string {
	return strconv.FormatFloat(float64(num), 'f', dig, 64)
}

func (ch *Time) SetTimebase(base float64) {
	ch.Base = Timebase(base)
}

func (ch Time) Secs() int {
	secs := ch.Duration / ch.Base.Float()
	return int(math.Round(secs))
}

func (ch Time) Float() float64 {
	return ch.Duration
}

func (ch Time) String() string {
	return ParseNumber(ch.Duration, 0)
}

func (ch Time) SecsString() string {
	return ParseNumber(ch.Duration/ch.Base.Float(), 3)
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
