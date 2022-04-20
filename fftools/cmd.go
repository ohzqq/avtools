package fftools

import (
	"os"
	"log"
	//"fmt"
	"strings"
	"os/exec"
	"strconv"
	"io/fs"
	"path/filepath"

	//"github.com/alessio/shellescape"
)

type FFmpegCmd struct {
	cmd *exec.Cmd
	Input []string
	profile bool
	padding bool
	args CmdArgs
}

func NewCmd() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.padding = false
	ff.cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) Profile(p string) *FFmpegCmd {
	ff.profile = true
	ff.args = Cfg.Profiles[p]
	ff.args.LogLevel(Cfg.Defaults.Verbosity)
	ff.args.Out(Cfg.Defaults.Output)
	ff.args.OverWrite(Cfg.Defaults.Overwrite)
	return ff
}

func (ff *FFmpegCmd) Args() *CmdArgs {
	if !ff.profile {
		ff.args = CmdArgs{}
		ff.args.VCodec("copy")
		ff.args.ACodec("copy")
	}
	return &ff.args
}

func (ff *FFmpegCmd) Join(ext string) *FFmpegCmd {
	files := find("." + ext)
	fileList := concatFile(files)
	//defer os.Remove(fileList.Name())
	ff.In(fileList.Name())
	pre := flagArgs{
		"f": "concat",
		"safe": "0",
	}
	ff.Args().PreInput(pre).Ext(ext)
	return ff
}

func concatFile(files []string) *os.File {
	file, err := os.CreateTemp("", "audiofiles")
	if err != nil { log.Fatal(err) }

	var cat strings.Builder
	for _, f := range files {
		abs, err := filepath.Abs(f)
		if err != nil { log.Fatal(err) }
		cat.WriteString("file ")
		cat.WriteString("'")
		cat.WriteString(abs)
		cat.WriteString("'")
		cat.WriteString("\n")
	}

	if _, err := file.WriteString(cat.String()); err != nil {
		log.Fatal(err)
	}

	return file
}

func find(ext string) []string {
	var files []string
	filepath.WalkDir(".", func(file string, dir fs.DirEntry, e error) error {
		if e != nil { return e }
		if filepath.Ext(dir.Name()) == ext {
			files = append(files, file)
		}
		return nil
	})
	return files
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
			if len(ff.Input) > 0 {
				for _, i := range ff.Input {
					ff.push("-i")
					ff.push(i)
				}
			}
		case "Meta":
			ff.Meta()
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
			if ff.Output() != "" {
				ff.cmd.Args = append(ff.cmd.Args, ff.Output())
			}
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

func (ff *FFmpegCmd) In(input string) {
	ff.Input = append(ff.Input, input)
}

func (ff *FFmpegCmd) Meta() {
	if ff.args.Metadata != "" {
		ff.push("-i")
		ff.push(ff.args.Metadata)
		ff.push("-map_metadata")
		ff.push(strconv.Itoa(len(ff.Input)))
	}
}

func (ff *FFmpegCmd) Pre() {
	for _, arg := range ff.args.Pre.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) Overwrite() {
	if Cfg.Defaults.Overwrite {
		ff.push("-y")
	}
}

func (ff *FFmpegCmd) Post() {
	for _, arg := range ff.args.Post.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) VideoCodec() {
	if ff.args.VideoCodec != "" {
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
	if ff.args.AudioCodec != "" {
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

func (ff *FFmpegCmd) Output() string {
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
		ext = "." + ff.args.Extension
	} else {
		ext = ".mkv"
	}

	if ff.args.Output == "" {
		return o + pad + ext
	} else {
		return ff.args.Output + pad + ext
	}

	return ""
}
