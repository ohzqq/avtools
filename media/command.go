package media

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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
	Meta     bool
	Cue      bool
	Cover    bool
	Chapters bool
}

type Files struct {
	Meta  string
	Cue   string
	Cover string
}

type UpdateCmd struct {
	*Media
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

func (cmd Command) Remove(input string) Cmd {
	m := New(input)

	f := m.Command()

	if cmd.Flags.Bool.Meta {
		f.Input.MapMetadata("-1")
	}

	if cmd.Flags.Bool.Chapters {
		f.Input.MapChapters("-1")
	}

	if cmd.Flags.Bool.Cover {
		f.Output.Set("vn", "")
	}

	f.Input.Set("y", "")
	name := m.Input.NewName().Prefix("updated-").Join()
	f.Output.Set("c", "copy")
	f.Output.Ext(m.Input.Ext).Name(name).Pad("")

	return f.Compile()
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
		ff := ExtractCover(m)
		cmds = append(cmds, ff)
	}

	return cmds
}

func Join(ext string, dir ...string) (Cmd, map[string]Cmd) {
	d := "."
	if len(dir) > 0 {
		d = dir[0]
	}

	path, err := filepath.Abs(d)
	if err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob(path + "/*" + ext)
	if err != nil {
		log.Fatal(err)
	}

	tmp, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer tmp.Close()

	var media []*Media
	for _, f := range files {
		media = append(media, New(f))
		str := "file '" + f + "'\n"
		if _, err := tmp.WriteString(str); err != nil {
			log.Fatal(err)
		}
	}

	formats := make(map[string]Cmd)
	tmpMedia := media[0]
	tmpMedia.SetChapters(GenerateChapters(media))
	formats["ini"] = tmpMedia.SaveMetaFmt("ini")
	formats["cue"] = tmpMedia.SaveMetaFmt("cue")

	cmd := ff.New()
	cmd.In(tmp.Name())
	cmd.Input.Set("f", "concat")
	cmd.Input.Set("safe", "0")
	cmd.Input.Set("y", "")

	if media[0].HasCover {
		cmd.Output.Set("vn", "")
	}

	base := filepath.Base(d)
	name := filepath.Join(path, base)
	cmd.Output.Set("c", "copy").Ext(ext).Name(name)

	return cmd.Compile(), formats
}

func GenerateChapters(media []*Media) []*avtools.Chapter {
	var chapters []*avtools.Chapter

	var start = []int64{0}
	var end []int64
	for idx, m := range media {
		d := m.GetTag("duration")
		dur := avtools.ParseStamp(d)
		e := start[idx] + dur.Milliseconds()
		end = append(end, e)
		if idx < len(media)-1 {
			s := end[idx]
			start = append(start, s)
		}
		ss := avtools.ParseStampDuration(start[idx], 1000)
		to := avtools.ParseStampDuration(end[idx], 1000)
		chapter := &avtools.Chapter{
			Start: avtools.Timestamp(ss),
			End:   avtools.Timestamp(to),
			Title: "Chapter " + strconv.Itoa(idx+1),
		}
		chapters = append(chapters, chapter)
	}

	return chapters
}
func ExtractCover(m *Media) Cmd {
	var stream Stream
	for _, s := range m.VideoStreams() {
		if s.IsCover {
			stream = s
		}
		break
	}
	cmd := ff.New()
	cmd.In(m.Input.Abs, ffmpeg.KwArgs{"y": ""})

	name := m.Input.NewName()
	n := name.Prefix("cover-").Join()
	cmd.Output.Name(n)
	cmd.Output.Pad("").Set("c", "copy").Set("an", "")

	switch stream.CodecName {
	case "mjpeg":
		cmd.Ext(".jpg")
	case "png":
		cmd.Ext(".png")
	}

	return cmd.Compile()
}

func (m Media) SaveMetaFmt(f string) Cmd {
	var cmd Cmd
	switch f {
	case "ini":
		name := m.Input.NewName()
		file := name.WithExt(".ini")
		file.Save(m.DumpIni())
		cmd = file
	case "ffmeta":
		cmd = m.DumpFFMeta()
	case "cue":
		if m.HasChapters() {
			name := m.Input.NewName()
			file := name.WithExt(".cue")
			file.Save(m.DumpCue())
			cmd = file
		}
	}
	return cmd
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

func CutChapter(media *Media, chapter *avtools.Chapter) Cmd {
	out := media.Input.NewName()

	title := chapter.Title
	if title == "" {
		title = fmt.Sprintf("-%s-%s", chapter.Start.Dur, chapter.End.Dur)
	}
	out.Suffix(title)
	//fmt.Printf("chapter %+V\n", chapter)

	cmd := media.Command()

	cmd.Input.Start(chapter.Start.String()).
		End(chapter.End.String())

	name := filepath.Join(out.Path, out.Name)
	cmd.Output.Set("c", "copy").
		Name(name).
		Pad("").
		Ext(media.Input.Ext)

	return cmd.Compile()
}

func (cmd Command) Update(input string) Cmd {
	m := UpdateCmd{
		Media: cmd.updateMeta(input),
	}

	return m
}

func (up UpdateCmd) Run() error {
	if up.MetaChanged {
		file := up.Input.NewName()
		file.Tmp(up.DumpIni())
		file.Run()
		tmp := file.file.Name()
		cmd := up.Command()
		cmd.Input.FFMeta(tmp)

		cmd.Output.Set("c", "copy")
		name := up.Input.NewName().Prefix("updated-").Join()
		cmd.Output.Ext(up.Input.Ext).Name(name).Pad("")

		c := cmd.Compile()
		//fmt.Println(c.String())

		err := c.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
