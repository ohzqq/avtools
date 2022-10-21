package tool

type SplitCmd struct {
	*Cmd
}

func Split() *SplitCmd {
	return &SplitCmd{Cmd: NewCmd()}
}

func (s *SplitCmd) Parse() *Cmd {
	for idx, _ := range s.Media.EachChapter() {
		cut := Cut()
		cut.ParseFlags(s.flag)
		ff := cut.Chap(idx + 1)
		s.Add(ff)
	}
	return s.Cmd
}
