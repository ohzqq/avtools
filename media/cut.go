package media

import (
	"fmt"
	"path/filepath"

	"github.com/ohzqq/avtools"
)

type CutCmd struct {
	*Media
	Chap *avtools.Chapter
}

func Cut(file string) CutCmd {
	cut := CutCmd{
		Media: New(file).Probe(),
		Chap:  &avtools.Chapter{},
	}
	return cut
}

func CutTime(file, start, end string) {
}

func (c *CutCmd) AllChapters() {
	if c.HasChapters() {
		for _, ch := range c.Chapters() {
			cmd := c.SetChapter(ch)
			cmd.Compile().Run()
		}
	}
}

func (c *CutCmd) Chapter(num int) *CutCmd {
	if c.HasChapters() {
		if num > 0 && num <= len(c.Chapters()) {
			c.Chap = c.Chapters()[num-1]
		}
	}
	return c
}

func (c CutCmd) SetChapter(ch *avtools.Chapter) CutCmd {
	c.Chap = ch
	return c
}

func (c *CutCmd) Start(ss string) *CutCmd {
	c.Chap.SS(ss)
	return c
}

func (c *CutCmd) End(to string) *CutCmd {
	c.Chap.To(to)
	return c
}

func (c CutCmd) Title() string {
	if c.Chap.Title != "" {
		return c.Chap.Title
	}
	t := fmt.Sprintf("-%s-%s", c.Chap.Start.Dur, c.Chap.End.Dur)
	return t
}

func (c CutCmd) Compile() Cmd {
	out := c.Input.NewName()
	out.Suffix(c.Title())
	name := filepath.Join(out.Path, out.Name)

	cmd := c.Media.Command()

	cmd.Input.Start(c.Chap.Start.String()).
		End(c.Chap.End.String())

	cmd.Output.Set("c", "copy").
		Name(name).
		Pad("").
		Ext(c.Media.Input.Ext)

	return cmd
}

func CutChapter(media *Media, chapter *avtools.Chapter) Cmd {
	out := media.Input.NewName()

	title := chapter.Title
	if title == "" {
		title = fmt.Sprintf("-%s-%s", chapter.Start.Dur, chapter.End.Dur)
	}
	out.Suffix(title)

	cmd := media.Command()

	cmd.Input.Start(chapter.Start.String()).
		End(chapter.End.String())

	name := filepath.Join(out.Path, out.Name)
	cmd.Output.Set("c", "copy").
		Name(name).
		Pad("").
		Ext(media.Input.Ext)

	return cmd
}
