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
	Start string
	End string
	ChapNo int
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
	input, err := filepath.Abs(c.Input)
	if err != nil {
		log.Fatal(err)
	}

	switch ext := filepath.Ext(input); ext {
	case ".m4a", ".m4b", ".mp3":
	case ".ini":
		c.ChapFlag = false
		c.MetaFlag = false
		c.CoverFlag = false
		c.MetaFile = input
	}

	switch {
	case c.ChapFlag, c.MetaFlag, c.CoverFlag:
		ff := NewFFprobeCmd().In(input)
		ff.verbose = true

		args := ff.Args()
		args.Verbosity("error").Format("json")
		if c.ChapFlag {
			args.Chapters()
		}
		if c.MetaFlag {
			args.Entries("format_tags")
		}

		ff.Run()
	case c.MetaFile != "":
		meta := ReadFFmetadata(input)
		fmt.Printf("%+V\n", len(meta.Chapters))
	}
}

func(c *Cmd) Split() error {
	chaps, err := c.getChapters()
	if err != nil {
		return err
	}

	for i, ch := range chaps {
		cmd := c.Cut(ch.StartToSeconds(), ch.EndToSeconds(), i)
		if c.Overwrite {
			cmd.Args().OverWrite()
		}
		cmd.Run()
	}
	return nil
}

func(c *Cmd) getChapters() ([]*Chapter, error) {
	switch {
	case c.CueFile != "":
		return ReadCueSheet(c.CueFile), nil
	case c.MetaFile != "":
		return ReadFFmetadata(c.MetaFile).Chapters, nil
	case c.Media.HasChapters():
		return c.Media.Meta.Chapters, nil
	default:
		return nil, fmt.Errorf("There are no chapters!")
	}
}

func(c *Cmd) Cut(ss, to string, no int) *FFmpegCmd {
	var (
		count = no + 1
		cmd = NewFFmpegCmd().In(c.Media)
		start = ss
		end = to
	)

	if c.ChapNo != 0 {
		chaps, err := c.getChapters()
		if err != nil {
			log.Fatal(err)
		}
		ch := chaps[c.ChapNo - 1]
		count = c.ChapNo
		start = ch.StartToSeconds()
		end = ch.EndToSeconds()
	}

	cmd.Args().Post("ss", start).Post("to", end).Out("tmp" + fmt.Sprintf("%06d", count)).Ext(c.Media.Ext)

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
