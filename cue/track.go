package cue

import "fmt"

type Track struct {
	title string
	start int
	end   int
}

func NewTrack() Track {
	return Track{}
}

func (t Track) Title() string {
	return t.title
}

func (t *Track) SetTitle(title string) *Track {
	t.title = title
	return t
}

func (t Track) Start() int {
	return t.start
}

func (t Track) Stamp() string {
	mm := t.start / 60
	ss := t.start % 60
	start := fmt.Sprintf("%02d:%02d:00", mm, ss)
	return start
}

func (t *Track) SetStart(secs int) *Track {
	t.start = secs
	return t
}

func (t Track) End() int {
	return t.end
}

func (t Track) Timebase() float64 {
	return 1
}
