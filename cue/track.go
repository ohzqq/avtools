package cue

import "time"

type Track struct {
	start time.Duration
	end   time.Duration
	title string
}

func (t Track) Start() time.Duration {
	return t.start
}
func (t Track) End() time.Duration {
	return t.end
}
func (t Track) Title() string {
	return t.title
}
