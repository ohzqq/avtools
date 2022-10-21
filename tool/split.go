package tool

type SplitCmd struct {
	*Cmd
}

func Split() *SplitCmd {
	return &SplitCmd{Cmd: NewCmd()}
}

func (s *SplitCmd) Parse() *Cmd {
	for idx, _ := range s.Media().Meta.Chapters.Chapters {
		cut := Cut()
		cut.SetFlags(s.flag)
		ff := cut.Chap(idx + 1)
		s.Add(ff)
	}
	return s.Cmd
}
