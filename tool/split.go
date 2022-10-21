package tool

import (
	"fmt"

	"github.com/ohzqq/avtools/file"
)

type SplitCmd struct {
	*Cmd
}

func Split() *SplitCmd {
	return &SplitCmd{Cmd: NewCmd()}
}

func (s *SplitCmd) Parse() *Cmd {
	for idx, ch := range s.Media().Meta.Chapters.Chapters {
		ss := ch.Start().SecsString()
		to := ch.End().SecsString()
		in := file.New(s.Media().Input.String())
		cut := Cut().Start(ss).End(to)
		cut.Input(in.Abs)
		ff := cut.FFmpegCmd()
		ff.Output(in.Pad(idx))
		fmt.Println(ff.String())
		//s.Add(cut)
	}
	return s.Cmd
}
