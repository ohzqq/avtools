package cue

import (
	"path/filepath"
	"strings"
)

type Sheet struct {
	file       string
	Audio      string
	Tracks     []Track
	titles     []string
	startTimes []int
}

func NewCueSheet(file string) *Sheet {
	return &Sheet{file: file}
}

func (s *Sheet) SetAudio(name string) *Sheet {
	s.Audio = name
	return s
}

func (s Sheet) File() string {
	return s.Audio
}

func (s Sheet) Ext() string {
	ext := filepath.Ext(s.Audio)
	ext = strings.TrimPrefix(ext, ".")
	return strings.ToUpper(ext)
}
