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
	cmd.FFmpegCmd = NewFFmpegCmd()

	return cmd
}

func (c *Cmd) Args() *CmdArgs {
	return c.args
}

func(c *Cmd) Extract() {
	switch {
	case c.ChapFlag:
		c.Media.FFmetaChapsToCue()
		break
	case c.CoverFlag:
		c.Args().ACodec("an").Out("cover").Ext(".jpg")
	case c.MetaFlag:
		c.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Out("ffmeta").Ext(".ini")
	}

	c.FFmpegCmd.In(c.Media).SetArgs(c.args).Profile(c.Profile).Run()
}

func(c *Cmd) Remove() {
	if c.ChapFlag {
		c.Args().Post("map_chapters", "-1")
	}

	if c.CoverFlag {
		c.Args().VCodec("vn")
	}

	if c.MetaFlag {
		c.Args().Post("map_metadata", "-1")
	}

	c.FFmpegCmd.In(c.Media).SetArgs(c.args).Profile(c.Profile).Run()
}

func(c *Cmd) Update() {
	if c.CoverFile != "" {
		switch codec := c.Media.AudioCodec(); codec {
		case "aac":
			c.addAacCover()
			break
		case "mp3":
			c.Args().Out("tmp-cover").Params(Mp3CoverArgs()).Cover(c.CoverFile)
		}
	}
	if c.MetaFile != "" {
		c.Args().Meta(c.MetaFile)
	}
	c.FFmpegCmd.In(c.Media).SetArgs(c.args).Profile(c.Profile).Run()
}

func(c *Cmd) addAacCover() {
	cpath, err := filepath.Abs(c.CoverFile)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("AtomicParsley", c.Media.Path, "--artwork", cpath, "--overWrite")
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

func(c *Cmd) Show() {
	switch {
	case c.ChapFlag, c.MetaFlag, c.CoverFlag:
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
}

func(c *Cmd) Split() {
	var chaps []*Chapter
	switch {
	case c.CueFile != "":
		chaps = ReadCueSheet(c.CueFile)
	case c.Media.HasChapters():
		chaps = c.Media.Meta.Chapters
	}

	for i, ch := range chaps {
		cmd := c.Media.Cut(ch.Start, ch.End, i)
		if c.Overwrite {
			cmd.Args().OverWrite()
		}
		cmd.Run()
	}
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

func(c *Cmd) Join() {
	var (
		ff = NewFFmpegCmd()
		files = find(c.Ext)
		cat strings.Builder
	)

	ff.Args().Pre(flagArgs{"f": "concat", "safe": "0"}).VCodec("vn").Ext(c.Ext)

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
	ff.Run()
}
