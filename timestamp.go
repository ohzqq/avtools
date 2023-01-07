package avtools

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Time struct {
	Duration float64
	Dur      time.Duration
	base     float64
	Tbase    Timebase
	Base
}

type Timebase float64

type Base struct {
	time float64
}

func (t Base) String() string {
	return "1/" + ParseNumber(t.time, 0)
}

func (t Base) Float() float64 {
	return t.time
}

func (t Timebase) String() string {
	return "1/" + ParseNumber(t.Float(), 0)
}

func (t Timebase) Float() float64 {
	return float64(t)
}

func Timestamp[N Number](t N, b ...N) Time {
	var base float64 = 1
	if len(b) > 0 {
		base = float64(b[0])
	} else {
		b = []N{1}
	}
	dur := float64(t)

	return Time{
		Duration: float64(dur),
		Tbase:    Timebase(base),
		Base: Base{
			time: base,
		},
	}
}

func ParseString(t string) Time {
	return Timestamp(ParseStamp(t).Milliseconds(), 1000)
}

func StringToFloat(t string) float64 {
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func ParseStamp(t string) time.Duration {
	var hh string
	var mm string
	var ss float64
	switch split := strings.Split(t, ":"); len(split) {
	case 3:
		hh = split[0] + "h"
		mm = split[1] + "m"
		ss = StringToFloat(split[2])
	case 2:
		mm = split[0] + "m"
		ss = StringToFloat(split[1])
	case 1:
		ss = StringToFloat(split[0])
	}
	ms := strconv.Itoa(int(ss*1000)) + "ms"
	stamp := fmt.Sprintf("%s%s%s", hh, mm, ms)
	dur, err := time.ParseDuration(stamp)
	if err != nil {
		log.Fatal(err)
	}

	println(stamp)
	return dur
}

func ParseNumber[N Number](num N, dig int) string {
	return strconv.FormatFloat(float64(num), 'f', dig, 64)
}

func (ch *Time) SetTimebase(base float64) {
	ch.Tbase = Timebase(base)
}

func (ch Time) Secs() int {
	secs := ch.Duration / ch.Tbase.Float()
	return int(math.Round(secs))
}

func (ch Time) Float() float64 {
	return ch.Duration
}

func (ch Time) String() string {
	return ParseNumber(ch.Duration, 0)
}

func (ch Time) SecsString() string {
	return ParseNumber(ch.Duration/ch.Tbase.Float(), 3)
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
