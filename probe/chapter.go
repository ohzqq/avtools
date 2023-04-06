package probe

import (
	"log"
	"time"

	"github.com/ohzqq/dur"
)

type Chapter struct {
	Base         string            `json:"time_base" ini:"TIMEBASE"`
	StartTime    string            `json:"start_time" ini:"START"`
	EndTime      string            `json:"end_time" ini:"END"`
	ChapterTitle string            `ini:"title"`
	CTags        map[string]string `json:"tags"`
}

func (c Chapter) Start() time.Duration {
	f, err := dur.Parse(c.StartTime)
	if err != nil {
		log.Fatal(err)
	}
	return f.Dur
}

func (c Chapter) End() time.Duration {
	f, err := dur.Parse(c.EndTime)
	if err != nil {
		log.Fatal(err)
	}
	return f.Dur
}

func (c Chapter) Title() string {
	if t, ok := c.Tags["title"]; ok {
		return t
	}
	return c.ChapterTitle
}

func (c Chapter) Tags() map[string]string {
	return c.CTags
}
