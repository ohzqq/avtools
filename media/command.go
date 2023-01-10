package media

import (
	"github.com/ohzqq/avtools"
)

type Cmd interface {
	Run() error
}

type Command struct {
	Flags
}

type Flags struct {
	Bool Bool
	File Files
}

type Bool struct {
	Meta  bool
	Cue   bool
	Cover bool
}

type Files struct {
	Meta  string
	Cue   string
	Cover string
}

func (cmd Command) updateMeta(input string) *Media {
	m := New(input)

	switch {
	case cmd.Flags.File.Meta != "":
		m.LoadIni(cmd.Flags.File.Meta)
		m.MetaChanged = true
	case cmd.Flags.File.Cue != "":
		m.LoadCue(cmd.Flags.File.Cue)
		m.MetaChanged = true
	}

	return m
}

func (cmd Command) Extract(input string) []Cmd {
	m := New(input)

	var cmds []Cmd

	if cmd.Flags.Bool.Cue {
		c := m.SaveMetaFmt("cue")
		cmds = append(cmds, c)
	}

	if cmd.Flags.Bool.Meta {
		c := m.SaveMetaFmt("ffmeta")
		cmds = append(cmds, c)
	}

	if cmd.Flags.Bool.Cover {
		ff := m.ExtractCover()
		cmds = append(cmds, ff)
	}

	return cmds
}

func (cmd Command) Update(input string) Cmd {
	m := UpdateCmd{
		Media: cmd.updateMeta(input),
	}

	return m
}

func (cmd Command) CutStamp(input, start, end string) Cmd {
	var (
		chapter = &avtools.Chapter{}
		media   = New(input)
		ss      = "0"
		to      = media.GetTag("duration")
	)

	if start != "" {
		ss = start
	}
	chapter.SS(ss)

	if end != "" {
		to = end
	}
	chapter.To(to)

	return CutChapter(media, chapter)
}

func (cmd Command) CutChapter(input string, num int) Cmd {
	media := New(input)
	chapter := media.GetChapter(num)
	return CutChapter(media, chapter)
}

func (cmd Command) Split(input string) []Cmd {
	media := cmd.updateMeta(input)

	var cmds []Cmd
	for _, chapter := range media.Chapters() {
		ch := CutChapter(media, chapter)
		cmds = append(cmds, ch)
	}

	return cmds
}
