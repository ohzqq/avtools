package avtools

import (
	"fmt"
	"os"
	"bytes"
	"os/exec"
	"log"
	"strings"
	"path/filepath"
)

type Cmd struct {
	Media *Media
	args *CmdArgs
	CliArgs *CmdArgs
	CmdArgs *CmdArgs
	Arg []string
	FFmpegCmd *FFmpegCmd
	FFprobeCmd *FFprobeCmd
	InputSlice []string
	Input string
	cwd string
}

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &Cmd{
		Media: NewMedia(),
		cwd: cwd,
	}
	//cmd.args = Cfg().GetProfile(cmd.Profile)
	//cmd.FFmpegCmd = NewFFmpegCmd()
	//fmt.Printf("%+v\n", cmd.Args())

	return cmd
}

func (c *Cmd) SetProfile(p string) *Cmd {
	c.args = Cfg().GetProfile(p)
	return c
}

func (c *Cmd) Args() *CmdArgs {
	return c.args
}

func(c *Cmd) ParseArguments() *Cmd {
	c.Media = NewMedia().Input(c.Input)

	switch {
	case c.CliArgs.Profile != "":
		c.CmdArgs.Profile = c.CliArgs.Profile
	default:
		c.CmdArgs.Profile = Cfg().DefaultProfile()
	}
	c.CmdArgs = Cfg().GetProfile(c.CmdArgs.Profile)

	// -log_level switch
	switch {
	case c.CmdArgs.LogLevel != "":
		c.Arg = append(c.Arg, "-loglevel", Cfg().defaults.LogLevel)
	}

	// -y switch
	switch {
	case c.CliArgs.Overwrite == true:
		c.Arg = append(c.Arg, "-y")
		break
	case Cfg().defaults.Overwrite == true:
		c.Arg = append(c.Arg, "-y")
	}

	// Anything before -i switches
	switch {
	case c.CliArgs.PreInput != nil:
		c.Arg = append(c.Arg, c.CliArgs.PreInput.Split()...)
	case c.CmdArgs.PreInput != nil:
		c.Arg = append(c.Arg, c.CmdArgs.PreInput.Split()...)
	}

	// -i switches

	// Anything after -i switches
	switch {
	case c.CliArgs.PostInput != nil:
		c.Arg = append(c.Arg, c.CliArgs.PostInput.Split()...)
	case c.CmdArgs.PostInput != nil:
		c.Arg = append(c.Arg, c.CmdArgs.PostInput.Split()...)
	}

	// -c:v
	//switch {
	//case c.CliArgs.VideoCodec != "":
	//  switch vc := ff.args.VideoCodec; vc {
	//  case "":
	//  case "none":
	//    fallthrough
	//  case "vn":
	//    ff.push("-vn")
	//  default:
	//    ff.push("-c:v")
	//    ff.push(ff.args.VideoCodec)
	//  }
	//  c.Arg = append(c.Arg, "-c:v", c.CliArgs.VideoCodec)
	//  break
	//case c.CmdArgs.VideoCodec != "":
	//  c.Arg = append(c.Arg, "-c:v", c.CmdArgs.VideoCodec)
	//}


	// any video codec params
	switch {
	case c.CliArgs.VideoParams != nil:
		c.Arg = append(c.Arg, c.CliArgs.VideoParams.Split()...)
	case c.CmdArgs.VideoParams != nil:
		c.Arg = append(c.Arg, c.CmdArgs.VideoParams.Split()...)
	}

	// any video filters
	if len(c.CmdArgs.VideoFilters) > 0 {
		c.Arg = append(c.Arg, "-vf", strings.Join(c.CmdArgs.VideoFilters, ","))
	}

	// filter complex
	if len(c.CmdArgs.FilterComplex) > 0 {
		c.Arg = append(c.Arg, "-vf", strings.Join(c.CmdArgs.FilterComplex, ","))
	}

	// -c:a
	switch {
	case c.CliArgs.AudioCodec != "":
		c.Arg = append(c.Arg, "-c:a", c.CliArgs.AudioCodec)
		break
	case c.CmdArgs.AudioCodec != "":
		c.Arg = append(c.Arg, "-c:a", c.CmdArgs.AudioCodec)
	}

	// any audio codec params
	switch {
	case c.CliArgs.AudioParams != nil:
		c.Arg = append(c.Arg, c.CliArgs.AudioParams.Split()...)
	case c.CmdArgs.AudioParams != nil:
		c.Arg = append(c.Arg, c.CmdArgs.AudioParams.Split()...)
	}

	// any audio filters
	if len(c.CmdArgs.AudioFilters) > 0 {
		c.Arg = append(c.Arg, "-af", strings.Join(c.CmdArgs.AudioFilters, ","))
	}

	switch {
	case c.CliArgs.Output != "":
		c.CmdArgs.Name = c.CliArgs.Output
	case c.CmdArgs.Output != "":
		c.CmdArgs.Name = Cfg().defaults.Output
	}

	//switch {
	//case c.CmdArgs.Padding != "":
	//  c.CmdArgs.Padding = proArgs.Padding
	//default:
	//  c.CmdArgs.Padding = Cfg().defaults.Padding
	//}

	switch {
	case c.CliArgs.Extension == "":
		c.CmdArgs.Extension = c.CliArgs.Extension
	}

	return c
}

func(c *Cmd) Extract() {
	switch {
	case c.Arg.ChapFlag:
		c.Media.FFmetaChapsToCue()
		break
	case c.Arg.CoverFlag:
		c.Args().ACodec("an").Out("cover").Ext(".jpg")
	case c.Meta:
		c.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Out("ffmeta").Ext(".ini")
	}

	if c.Profile != "" {
		c.FFmpegCmd.Profile(c.Profile)
	}

	c.FFmpegCmd.In(c.Media).SetArgs(c.args).Run()
}

func(c *Cmd) Remove() {
	if c.Arg.ChapFlag {
		c.Args().Post("map_chapters", "-1")
	}

	if c.Arg.CoverFlag {
		c.Args().VCodec("vn")
	}

	if c.Arg.MetaFlag {
		c.Args().Post("map_metadata", "-1")
	}

	if c.Profile != "" {
		c.FFmpegCmd.Profile(c.Arg.Profile)
	}

	c.FFmpegCmd.In(c.Media).SetArgs(c.args).Run()
}

func(c *Cmd) Update() {
	if c.Arg.CoverFile != "" {
		switch codec := c.Media.AudioCodec(); codec {
		case "aac":
			c.addAacCover()
			break
		case "mp3":
			c.Args().Out("tmp-cover").Params(Mp3CoverArgs()).Cover(c.CoverFile)
		}
	}
	if c.Arg.MetaFile != "" {
		//c.Args().Meta(c.Arg.MetaFile)
	}

	if c.Arg.Profile != "" {
		c.FFmpegCmd.Profile(c.Arg.Profile)
	}

	c.FFmpegCmd.In(c.Media).SetArgs(c.args).Run()
}

func(c *Cmd) addAacCover() {
	cpath, err := filepath.Abs(c.Arg.CoverFile)
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
		c.Arg.ChapFlag = false
		c.Arg.MetaFlag = false
		c.Arg.CoverFlag = false
		c.Arg.MetaFile = input
	}

	switch {
	case c.Arg.ChapFlag, c.Arg.MetaFlag, c.Arg.CoverFlag:
		ff := NewFFprobeCmd().In(input)
		ff.verbose = true

		args := ff.Args()
		args.Verbosity("error").Format("json")
		if c.Arg.ChapFlag {
			args.Chapters()
		}
		if c.Arg.MetaFlag {
			args.Entries("format_tags")
		}

		ff.Run()
	case c.Arg.MetaFile != "":
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
	case c.Arg.CueFile != "":
		return ReadCueSheet(c.Arg.CueFile), nil
	case c.Arg.MetaFile != "":
		return ReadFFmetadata(c.Arg.MetaFile).Chapters, nil
	case c.Media.HasChapters():
		return c.Media.Meta.Chapters, nil
	default:
		return nil, fmt.Errorf("There are no chapters!")
	}
}

func(c *Cmd) ffmpeg() *FFmpegCmd {
	if c.Arg.Profile != "" {
		c.SetProfile(c.Arg.Profile)
	}

	//if c.Args() == "" {
	//  c.Args().Out(Cfg().defaults)
	//}

	if c.Args().Padding == "" {
		c.Args().Pad(Cfg().defaults.Padding)
	}
	//fmt.Printf("%+v\n", c.Args())
	return NewFFmpegCmd().In(c.Media).SetArgs(c.args)
}

func(c *Cmd) Cut(ss, to string, no int) *FFmpegCmd {
	var (
		count = no + 1
		start = ss
		end = to
	)

	if c.Arg.ChapNo != 0 {
		chaps, err := c.getChapters()
		if err != nil {
			log.Fatal(err)
		}
		ch := chaps[c.Arg.ChapNo - 1]
		count = c.Arg.ChapNo
		start = ch.StartToSeconds()
		end = ch.EndToSeconds()
	}

	c.Args().
		Pre("ss", start).
		Pre("to", end).
		Num(count).
		Out("tmp")

	return c.ffmpeg()
}

func(c *Cmd) Join() {
	var (
		ff = NewFFmpegCmd()
		files = find(c.Arg.Extension)
		cat strings.Builder
	)

	if c.Arg.Profile != "" {
		ff.Profile(c.Arg.Profile)
	}

	ff.Args().
		Pre("f", "concat").
		Pre("safe", "0").
		VCodec("vn").Ext(c.Arg.Extension)

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
