package avtools

import "github.com/ohzqq/avtools/timestamp"

type Chapter struct {
	start timestamp.Time
	end   timestamp.Time
	base  timestamp.Timebase
	Title string
}
