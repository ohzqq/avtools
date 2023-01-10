package media

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
		Media: New(input),
	}

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
