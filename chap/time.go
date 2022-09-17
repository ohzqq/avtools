package chap

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

func ParseStr(t string) ChTime {
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		log.Fatal(err)
	}
	return ChTime{time: i}
}

type ChTime struct {
	time float64
	base float64
}

func NewChapterTime[N Number](num N) ChTime {
	return ChTime{
		time: float64(num),
		base: 1,
	}
}

func ParseNum[N Number](num N) string {
	return strconv.FormatFloat(float64(num), 'f', 0, 64)
}

func (ch ChTime) Secs() int {
	secs := ch.time / ch.base
	return int(math.Round(secs))
}

func (c ChTime) HHMMSS() string {
	secs := c.Secs()
	hh := secs / 3600
	mm := secs % 3600 / 60
	ss := secs % 3600 % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

func (c ChTime) MMSS() string {
	secs := c.Secs()
	mm := secs / 60
	ss := secs % 60
	return fmt.Sprintf("%02d:%02d", mm, ss)
}
