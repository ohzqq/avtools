package cue

import (
	"path/filepath"
	"strings"
)

type CueSheet struct {
	file       string
	Audio      string
	Tracks     []Track
	titles     []string
	startTimes []int
}

func NewCueSheet(file string) *CueSheet {
	return &CueSheet{file: file}
}

func (c *CueSheet) SetAudio(name string) *CueSheet {
	c.Audio = name
	return c
}

func (c CueSheet) File() string {
	return c.Audio
}

func (c CueSheet) Ext() string {
	ext := filepath.Ext(c.Audio)
	ext = strings.TrimPrefix(ext, ".")
	return strings.ToUpper(ext)
}
