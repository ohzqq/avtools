package yygif

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools/ff"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

type Meta map[string]Scene

type Scene map[string]Clip

type Clip struct {
	Name  string
	Video string `yaml:"video"`
	Start string `yaml:"s"`
	End   string `yaml:"e"`
	Crop  string `yaml:"crop"`
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
	return ffmpeg.KwArgs{"ss": c.Start, "to": c.End}
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
			scenes[k] = clip
		}
		meta[scene] = scenes
	}

	return meta
}
