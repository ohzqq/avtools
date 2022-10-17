package chap

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
	time float64
	base float64
}

func NewChapterTime[N Number](num N) Time {
	return Time{
		time: float64(num),
		base: 1,
	}
}

func ParseStr(t string) Time {
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		log.Fatal(err)
	}
	return Time{time: i}
}

func ParseNum[N Number](num N) string {
	return strconv.FormatFloat(float64(num), 'f', 0, 64)
}

func (ch *Time) SetTimebase(base float64) {
	ch.base = base
}

func (ch Time) Secs() int {
	secs := ch.time / ch.base
	return int(math.Round(secs))
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
