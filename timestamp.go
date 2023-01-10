package avtools

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Time struct {
	hh  int
	mm  int
	ss  float64
	Dur time.Duration
}

func Timestamp(d time.Duration) Time {
	stamp := Time{Dur: d}

	t := strings.SplitAfter(d.String(), "m")
	ss := strings.TrimSuffix(t[len(t)-1], "s")
	stamp.ss = StringToFloat(ss)

	t = lo.DropRight(t, 1)

	if len(t) > 0 {
		t = strings.SplitAfter(t[0], "h")
		mm := strings.TrimSuffix(t[len(t)-1], "m")
		stamp.mm, _ = strconv.Atoi(mm)
		t = lo.DropRight(t, 1)

		if len(t) > 0 {
			hh := strings.TrimSuffix(t[len(t)-1], "h")
			stamp.hh, _ = strconv.Atoi(hh)
		}
	}

	return stamp
}

func ParseString(t string) Time {
	return Timestamp(ParseStamp(t))
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
	var ss string
	switch split := strings.Split(t, ":"); len(split) {
	case 3:
		hh = split[0] + "h"
		mm = split[1] + "m"
		ss = split[2] + "s"
	case 2:
		mm = split[0] + "m"
		ss = split[1] + "s"
	case 1:
		ss = split[0] + "s"
	}
	stamp := fmt.Sprintf("%s%s%s", hh, mm, ss)

	return ParseDuration(stamp)
}

func ParseTimeAndBase(t, b string) time.Duration {
	base, err := strconv.Atoi(strings.TrimPrefix(b, "1/"))
	if err != nil {
		log.Fatal(err)
	}

	dur, err := strconv.Atoi(t)
	if err != nil {
		log.Fatal(err)
	}

	secs := dur / base
	ms := secs * 1000
	d := strconv.Itoa(ms) + "ms"

	return ParseDuration(d)
}

func ParseStampDuration[N Number](t, b N) time.Duration {
	secs := float64(t) / float64(b)
	ms := secs * 1000
	d := strconv.Itoa(int(ms)) + "ms"
	return ParseDuration(d)
}

func ParseDuration(d string) time.Duration {
	dur, err := time.ParseDuration(d)
	if err != nil {
		log.Fatal(err)
	}
	return dur
}

func ParseNumber[N Number](num N, dig int) string {
	return strconv.FormatFloat(float64(num), 'f', dig, 64)
}

func (ch Time) secs() int {
	ss := int(math.Round(ch.ss))
	return ss
}

func (c Time) String() string {
	return fmt.Sprintf("%02d:%02d:%06.3f", c.hh, c.mm, c.ss)
}

func (c Time) MS() string {
	ms := strconv.Itoa(int(c.Dur.Milliseconds()))
	return ms
}

func (ch Time) Min() int {
	min := ch.Dur.Minutes()
	return int(min)
}

func (c Time) HHMMSS() string {
	return fmt.Sprintf("%02d:%02d:%02d", c.hh, c.mm, c.secs())
}

func (c Time) MMSS() string {
	return fmt.Sprintf("%02d:%02d", c.Min(), c.secs())
}
