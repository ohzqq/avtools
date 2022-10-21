package tool

import (
	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/file"
)

type CutCmd struct {
	*Cmd
	ss string
	to string
}

func Cut() *CutCmd {
	return &CutCmd{
		Cmd: NewCmd(),
	}
}

func (c CutCmd) Chap(no int) *ffmpeg.Cmd {
	var (
		num   = no - 1
		chaps = c.Media().Meta.Chapters.Chapters
		ff    *ffmpeg.Cmd
	)

	if num < len(chaps) {
		ch := chaps[num]
		ss := ch.Start().SecsString()
		to := ch.End().SecsString()
		in := file.New(c.Media().Input.String())
		c.Start(ss).End(to)
		ff = c.FFmpegCmd()
		ff.Output(in.Pad(no))
	}

	return ff
}

func (c *CutCmd) Start(ss string) *CutCmd {
	c.ss = ss
	return c
}

func (c *CutCmd) End(to string) *CutCmd {
	c.to = to
	return c
}

func (c *CutCmd) Parse() *Cmd {
	if c.flag.Args.HasStart() {
		c.ss = c.flag.Args.Start
	}

	if c.flag.Args.HasEnd() {
		c.to = c.flag.Args.End
	}

	ff := c.FFmpegCmd()

	c.Add(ff)

	return c.Cmd
}

func (c *CutCmd) FFmpegCmd() *ffmpeg.Cmd {
	ff := c.FFmpeg()
	ff.SS(c.ss)
	ff.To(c.to)

	return ff
}
