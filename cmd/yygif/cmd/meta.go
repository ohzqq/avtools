package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "A brief description of your command",
	//Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if MetaExists("metadata-default.yml") {
			gifMeta := ReadMeta("metadata-default.yml")
			ini := gifMeta.DumpIni()
			file, err := os.Create("gif-meta.ini")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			_, err = file.Write(ini)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

type Meta map[string]Scene

type Scene map[string]Clip

type Clip struct {
	Name      string
	Video     string `yaml:"video"`
	StartTime string `yaml:"s"`
	EndTime   string `yaml:"e"`
	Crop      string `yaml:"crop"`
}

func (m Meta) Chapters() []*avtools.Chapter {
	var chapters []*avtools.Chapter

	count := 1
	for _, scene := range m {
		for _, clip := range scene {
			clip.Name = fmt.Sprintf("Gif%03d", count)
			ch := avtools.NewChapter(clip)
			ch.Tags["crop"] = clip.Crop
			count++
			chapters = append(chapters, ch)
		}
	}
	return chapters
}

func (m Meta) Streams() []map[string]string {
	return []map[string]string{}
}

func (m Meta) Tags() map[string]string {
	clip := m.GetClip("1", "1")
	return map[string]string{
		"title": clip.Input.File,
		//"title": media.NewFile(clip.Input.File).Abs,
	}
}

func (meta Meta) DumpIni() []byte {
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	ffmeta := ini.Empty(opts)

	for k, v := range meta.Tags() {
		_, err := ffmeta.Section("").NewKey(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, chapter := range meta.Chapters() {
		sec, err := ffmeta.NewSection("CHAPTER")
		if err != nil {
			log.Fatal(err)
		}
		sec.NewKey("START", chapter.Start.String())
		sec.NewKey("END", chapter.End.String())
		sec.NewKey("title", chapter.Title)
		for k, v := range chapter.Tags {
			sec.NewKey(k, v)
		}
	}

	var buf bytes.Buffer
	_, err := buf.WriteString(";FFMETADATA1\n")
	_, err = ffmeta.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (c Clip) Chap() *avtools.Chapter {
	ch := avtools.NewChapter(c)
	ch.Tags["crop"] = c.Crop
	return ch
}

func (c Clip) Start() time.Duration {
	return avtools.ParseStamp(c.StartTime)
}

func (c Clip) End() time.Duration {
	return avtools.ParseStamp(c.EndTime)
}

func (c Clip) Title() string {
	return c.Name
}

func (m Meta) MkGifs() []*ff.Cmd {
	var cmds []*ff.Cmd
	for s, scene := range m {
		for c, _ := range scene {
			s = strings.TrimPrefix(s, "scene")
			c = strings.TrimPrefix(c, "clip")
			gif := m.GetClip(s, c)
			cmds = append(cmds, gif)
		}
	}
	return cmds
}

func (m Meta) GetClip(s, c string) *ff.Cmd {
	sNum, _ := strconv.Atoi(s)
	cNum, _ := strconv.Atoi(c)
	scene := fmt.Sprintf("scene%03d", sNum)
	cKey := fmt.Sprintf("clip%03d", cNum)
	name := fmt.Sprintf("%s-clip", scene)

	clip := m[scene][cKey]

	cmd := clip.Compile()

	cmd.Output.Name(name).Num(cNum)

	return cmd.Compile()
}

func (c Clip) Compile() ff.Cmd {
	cmd := ff.New("gif")
	cmd.In(c.Video, c.Input())
	cmd.Filters = c.Filters()
	cmd.Output = c.Output()
	return cmd
}

func (c Clip) Input() ffmpeg.KwArgs {
	return ffmpeg.KwArgs{"ss": c.StartTime, "to": c.EndTime}
}

func (c Clip) Filters() ff.Filters {
	filters := ff.GetProfile("gif").Filters
	if c.Crop != "" {
		crop := strings.Split(c.Crop, ":")
		filters.Set("crop", crop...)
	}
	return filters
}

func (c Clip) Output() ff.Output {
	out := ff.GetProfile("gif").Output
	if _, ok := out.Args["c"]; ok {
		out.Del("c")
	}
	return out
}

func ReadMeta(name string) Meta {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	var yml map[string]yaml.Node
	err = yaml.Unmarshal(data, &yml)
	if err != nil {
		log.Fatal(err)
	}

	meta := make(Meta)
	for scene, c := range yml {
		var clips map[string]yaml.Node
		err := c.Decode(&clips)
		if err != nil {
			log.Fatal(err)
		}

		scenes := make(map[string]Clip)

		var video string
		v := clips["video"]
		err = v.Decode(&video)
		if err != nil {
			log.Fatal(err)
		}
		delete(clips, "video")

		for k, v := range clips {
			var clip Clip
			err := v.Decode(&clip)
			if err != nil {
				log.Fatal(err)
			}
			clip.Video = video
			clip.Name = scene + k
			scenes[k] = clip
		}
		meta[scene] = scenes
	}

	return meta
}
func init() {
	rootCmd.AddCommand(metaCmd)
}
