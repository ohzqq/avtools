package ffmeta

import (
	"strconv"
	"strings"
)

type Chapter struct {
	Base         string            `json:"time_base",ini:"timebase"`
	StartTime    int               `json:"start",ini:"start"`
	EndTime      int               `json:"end",ini:"end"`
	ChapterTitle string            `ini:"title"`
	Tags         map[string]string `json:"tags"`
}

func (c Chapter) Title() string {
	return c.ChapterTitle
}

func (c Chapter) Start() int {
	return c.StartTime
}

func (c Chapter) End() int {
	return c.EndTime
}

func (c Chapter) Timebase() float64 {
	if tb := c.Base; tb != "" {
		c.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.ParseFloat(c.Base, 64)
	return baseFloat
}
