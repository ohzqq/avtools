package fftools

import (
	"os"
	"log"
	"fmt"
	"bytes"
	//"strings"
	"os/exec"
	"strconv"
	//"io/fs"
	//"path/filepath"

	//"github.com/alessio/shellescape"
)

type FFmpegCmd struct {
	cmd *exec.Cmd
	Input []*Media
	profile bool
	padding bool
	tmpFile *os.File
	args CmdArgs
}

func NewCmd() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.padding = false
	ff.cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) In(input *Media) *FFmpegCmd {
	ff.Input = append(ff.Input, input)
	return ff
}

func (ff *FFmpegCmd) Profile(p string) *FFmpegCmd {
	ff.profile = true
	if ff.args.PreInput == nil {
		ff.args.Pre(Cfg.Profiles[p].PreInput)
	}
	if ff.args.PostInput == nil {
		ff.args.Post(Cfg.Profiles[p].PostInput)
	}
	if ff.args.VideoParams == nil {
		ff.args.VParams(Cfg.Profiles[p].VideoParams)
	}
	if ff.args.VideoCodec == "" {
		ff.args.VCodec(Cfg.Profiles[p].VideoCodec)
	}
	if ff.args.VideoFilters == "" {
		ff.args.VFilters(Cfg.Profiles[p].VideoFilters)
	}
	if ff.args.AudioParams == nil {
		ff.args.AParams(Cfg.Profiles[p].AudioParams)
	}
	if ff.args.AudioCodec == "" {
		ff.args.ACodec(Cfg.Profiles[p].AudioCodec)
	}
	if ff.args.AudioFilters == "" {
		ff.args.AFilters(Cfg.Profiles[p].AudioFilters)
	}
	if ff.args.FilterComplex == "" {
		ff.args.Filter(Cfg.Profiles[p].FilterComplex)
	}
	if ff.args.Verbosity == "" {
		ff.args.LogLevel(Cfg.Defaults.Verbosity)
	}
	if ff.args.Output == "" {
		ff.args.Out(Cfg.Defaults.Output)
	}
	if ff.args.Overwrite == false {
		ff.args.OverWrite(Cfg.Defaults.Overwrite)
	}
	return ff
}

func (ff *FFmpegCmd) Meta() *MediaMeta {
	return ff.Input[0].ReadMeta()
}

func (ff *FFmpegCmd) Args() *CmdArgs {
	if !ff.profile {
		ff.args = CmdArgs{}
		ff.args.VCodec("copy")
		ff.args.ACodec("copy")
	}
	return &ff.args
}

func (ff *FFmpegCmd) GetChapters() *Chapters {
	var (
		meta *MediaMeta
	)
	input := ff.Input[0]
	if input.Cue != "" {
		meta = ReadCueSheet(input.Cue)
	} else if input.FFmeta != "" {
		meta = ReadFFmetadata(input.FFmeta)
	} else {
		meta = input.ReadMeta()
	}
	return meta.Chapters
}

func (ff *FFmpegCmd) Run() {
	if ff.tmpFile != nil {
		defer os.Remove(ff.tmpFile.Name())
	}

	cmd := ff.Cmd()

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	fmt.Printf("%v\n", cmd.String())

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		fmt.Printf("%v\n", stderr.String())
	}
	fmt.Printf("%v\n", stdout.String())
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	argOrder := []string{"Verbosity", "Overwrite", "Pre", "Input", "Meta", "Post", "VideoCodec", "VideoParams", "VideoFilters", "FilterComplex", "AudioCodec", "AudioParams", "AudioFilters", "Output", "Ext"}

	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.Verbosity()
		case "Overwrite":
			ff.Overwrite()
		case "Pre":
			ff.Pre()
		case "Input":
			ff.pushInput()
		case "Meta":
			ff.Metadata()
		case "Post":
			ff.Post()
		case "VideoCodec":
			ff.VideoCodec()
		case "VideoParams":
			ff.VideoParams()
		case "VideoFilters":
			ff.VideoFilters()
		case "FilterComplex":
			ff.FilterComplex()
		case "AudioCodec":
			ff.AudioCodec()
		case "AudioParams":
			ff.AudioParams()
		case "AudioFilters":
			ff.AudioFilters()
		case "Output":
			ff.Output()
		}
	}
	return ff.cmd
}

func (ff *FFmpegCmd) push(arg string) {
	ff.cmd.Args = append(ff.cmd.Args, arg)
}

func (ff *FFmpegCmd) Verbosity() {
	if Cfg.Defaults.Verbosity != "" {
		ff.push("-loglevel")
		ff.push(Cfg.Defaults.Verbosity)
	}
}

func (ff *FFmpegCmd) pushInput() {
	if len(ff.Input) > 0 {
		for _, i := range ff.Input {
			ff.push("-i")
			ff.push(i.Path)
		}
	} else {
		log.Fatal("No input specified")
	}
}

func (ff *FFmpegCmd) Metadata() {
	if ff.args.Metadata != "" {
		ff.push("-i")
		ff.push(ff.args.Metadata)
		ff.push("-map_metadata")
		ff.push(strconv.Itoa(len(ff.Input)))
	}
}

func (ff *FFmpegCmd) Pre() {
	for _, arg := range ff.args.PreInput.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) Overwrite() {
	if Cfg.Defaults.Overwrite {
		ff.push("-y")
	}
}

func (ff *FFmpegCmd) Post() {
	for _, arg := range ff.args.PostInput.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) VideoCodec() {
	switch vc := ff.args.VideoCodec; vc {
	case "none":
	case "":
	default:
		ff.push("-c:v")
		ff.push(ff.args.VideoCodec)
	}
}

func (ff *FFmpegCmd) VideoParams() {
	for _, arg := range ff.args.VideoParams.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) VideoFilters() {
	if ff.args.VideoFilters != "" {
		ff.push("-vf")
		ff.push(ff.args.VideoFilters)
	}
}

func (ff *FFmpegCmd) FilterComplex() {
	if ff.args.FilterComplex != "" {
		ff.push("-filter")
		ff.push(ff.args.FilterComplex)
	}
}

func (ff *FFmpegCmd) AudioCodec() {
	switch ac := ff.args.AudioCodec; ac {
	case "none":
	case "":
	default:
		ff.push("-c:a")
		ff.push(ff.args.AudioCodec)
	}
}

func (ff *FFmpegCmd) AudioParams() {
	for _, arg := range ff.args.AudioParams.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) AudioFilters() {
	if ff.args.AudioFilters != "" {
		ff.push("-af")
		ff.push(ff.args.AudioFilters)
	}
}

func (ff *FFmpegCmd) Output() {
	var o string
	var pad string
	var ext string
	if Cfg.Defaults.Output != "" {
		o = Cfg.Defaults.Output
		if ff.padding {
			pad = "%06d"
		} else {
			pad = ""
		}
	}
	if ff.args.Extension != "" {
		ext = ff.args.Extension
	} else {
		ext = ".mkv"
	}

	if ff.args.Output == "" {
		ff.push(o + pad + ext)
	} else {
		ff.push(ff.args.Output + pad + ext)
	}
}
