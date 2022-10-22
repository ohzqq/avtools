package cue

import (
	"strings"

	"github.com/ohzqq/avtools/file"
)

type Sheet struct {
	file       file.File
	Audio      file.File
	Tracks     []Track
	titles     []string
	startTimes []int
}

func NewCueSheet(f string) *Sheet {
	return &Sheet{file: file.New(f)}
}

func (s *Sheet) SetAudio(name string) *Sheet {
	s.Audio = file.New(name)
	return s
}

func (s Sheet) File() string {
	return s.Audio.Abs
}

func (s Sheet) Ext() string {
	ext := strings.TrimPrefix(s.Audio.Ext, ".")
	return strings.ToUpper(ext)
}
