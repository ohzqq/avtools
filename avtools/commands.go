package avtools

import (
	"fmt"
	//"io/fs"
	"os"
	"bytes"
	"os/exec"
	"log"
	"strings"
	//"strconv"
	"path/filepath"
)

type Cmd struct {
	Media *Media
	args *CmdArgs
	FFmpegCmd *FFmpegCmd
	FFprobeCmd *FFprobeCmd
	InputSlice []string
	Profile string
	Output string
	Input string
	Ext string
	CoverFile string
	MetaFile string
	CueFile string
	Verbose bool
	Overwrite bool
	CoverFlag bool
	MetaFlag bool
	CueFlag bool
	ChapFlag bool
	cwd string
	action string
	isFFmpeg bool
	isFFprobe bool
}

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &Cmd{
		Profile: Cfg().DefaultProfile(),
		Media: NewMedia(),
		cwd: cwd,
		args: NewArgs(),
	}
	cmd.FFmpegCmd = NewFFmpegCmd().Profile(cmd.Profile)
	cmd.FFprobeCmd = NewFFprobeCmd()

	return cmd
}

func(c *Cmd) Action(a string) *Cmd {
	c.action = a
	switch a {
	case "show":
		c.isFFprobe = true
	case "update":
		c.isFFmpeg = true
		switch {
		case c.CoverFile != "":
			switch codec := c.Media.AudioCodec(); codec {
			case "aac":
				c.isFFmpeg = false
				addAacCover(c.Media.Path, c.CoverFile)
				break
			case "mp3":
				c.Args().Out("tmp-cover").Params(Mp3CoverArgs()).Cover(c.CoverFile)
			}
		case c.MetaFile != "":
			c.Args().Meta(c.MetaFile)
			//m.AddFFmeta(c.MetaFile)
		}
	case "extract":
		c.isFFmpeg = true
		switch {
		case c.ChapFlag:
			c.Media.FFmetaChapsToCue()
			break
		case c.CoverFlag:
			c.Args().ACodec("an").Out("cover").Ext(".jpg")
		case c.MetaFlag:
			c.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Ext(".ini").Out("ffmeta")
		}
	case "remove":
		c.isFFmpeg = true
		if c.ChapFlag {
			c.Args().Post("map_chapters", "-1")
		}

		if c.CoverFlag {
			c.Args().VCodec("vn")
		}

		if c.MetaFlag {
			c.Args().Post("map_metadata", "-1")
		}
	}
	return c
}

func (c *Cmd) Args() *CmdArgs {
	return c.args
}

func(c *Cmd) Run() {
	fmt.Printf("%+V\n", c.args)
	switch {
	case c.isFFmpeg:
		cmd := NewFFmpegCmd().In(c.Media).SetArgs(c.args).Profile(c.Profile)
		cmd.Run()
	case c.isFFprobe:
	}
}

func addAacCover(path, cover string) {
	cpath, err := filepath.Abs(cover)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("AtomicParsley", path, "--artwork", cpath, "--overWrite")
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		log.Printf("%v finished with error: %v", cmd.String(), err)
		fmt.Printf("%v\n", stderr.String())
	}
}

func(m *Media) AddFFmeta(meta string) *FFmpegCmd {
	path, err := filepath.Abs(meta)
	if err != nil {
		log.Fatal(err)
	}
	cmd := NewFFmpegCmd().In(m)
	cmd.Args().Meta(path)
	return cmd
}

func(c *Cmd) Show() {
	c.FFprobeCmd.In(c.Input)
	c.FFprobeCmd.verbose = true

	args := c.FFprobeCmd.Args()
	args.Verbosity("error").Format("json")
	if c.ChapFlag {
		args.Chapters()
	}
	if c.MetaFlag {
		args.Entries("format_tags")
	}

	c.FFprobeCmd.Run()
}

func(m *Media) Remove(chaps, cover, meta bool) {
	cmd := NewFFmpegCmd().In(m)

	if chaps {
		cmd.Args().Post("map_chapters", "-1")
	}

	if cover {
		cmd.Args().VCodec("vn")
	}

	if meta {
		cmd.Args().Post("map_metadata", "-1")
	}

	cmd.Run()
	//fmt.Println("Success!")
}

func(m *Media) Extract(chaps, cover, meta bool) {
	cmd := NewFFmpegCmd().In(m)

	switch {
	case chaps:
		m.FFmetaChapsToCue()
	case cover:
		cmd.Args().ACodec("an").Out("cover").Ext(".jpg")
		cmd.Run()
	case meta:
		cmd.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Ext(".ini").Out("ffmeta")
		cmd.Run()
	}
}

func(m *Media) Split(cue string) {
	var chaps []*Chapter
	switch {
	case cue != "":
		chaps = ReadCueSheet(cue)
	case m.HasChapters():
		chaps = m.Meta.Chapters
	}

	for i, ch := range chaps {
		cmd := m.Cut(ch.Start, ch.End, i)
		if m.Overwrite {
			cmd.Args().OverWrite()
		}
		cmd.Run()
	}

	fmt.Println("Success!")
}

func(m *Media) Cut(ss, to string, no int) *FFmpegCmd {
	var (
		count = fmt.Sprintf("%06d", no + 1)
		cmd = NewFFmpegCmd().In(m)
	)

	if ss != "" {
		cmd.Args().Post("ss", ss)
	}

	if to != "" {
		cmd.Args().Post("to", to)
	}

	cmd.Args().Out("tmp" + count).Ext(m.Ext)

	return cmd
}

func Join(ext string) *FFmpegCmd {
	var (
		ff = NewFFmpegCmd()
		files = find(ext)
		cat strings.Builder
	)

	ff.Args().Pre(flagArgs{"f": "concat", "safe": "0"}).VCodec("vn").Ext(ext)

	file, err := os.CreateTemp("", "audiofiles")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		cat.WriteString("file ")
		cat.WriteString("'")
		cat.WriteString(f)
		cat.WriteString("'")
		cat.WriteString("\n")
	}

	if _, err := file.WriteString(cat.String()); err != nil {
		log.Fatal(err)
	}

	ff.tmpFile = file

	ff.In(NewMedia().Input(ff.tmpFile.Name()))

	return ff
}


