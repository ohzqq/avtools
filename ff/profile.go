package ff

import (
	"log"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

//var profiles = make(map[string]profile)

var profiles = map[string]profile{
	"quiet": profile{
		In: ffmpeg.KwArgs{
			"loglevel":    "error",
			"hide_banner": "",
		},
		Out: ffmpeg.KwArgs{
			"padding": "%03d",
			"name":    "tmp",
			"num":     1,
		},
	},
	"defaultAudio": profile{
		Out: ffmpeg.KwArgs{
			"ext": ".mka",
		},
	},
	"defaultVideo": profile{
		Out: ffmpeg.KwArgs{
			"ext": ".mkv",
		},
	},
	"stream": profile{
		Out: ffmpeg.KwArgs{
			"c:a": "copy",
			"c:v": "copy",
		},
	},
}

type profile struct {
	Filters Filters       `yaml:"filters"`
	In      ffmpeg.KwArgs `yaml:"input"`
	Out     ffmpeg.KwArgs `yaml:"output"`
}

func MergeProfiles(pros ...string) profile {
	var in []ffmpeg.KwArgs
	var out []ffmpeg.KwArgs
	var filters []Filters
	for _, p := range pros {
		if pro, ok := profiles[p]; ok {
			in = append(in, pro.In)
			out = append(out, pro.Out)
			filters = append(filters, pro.Filters)
		}
	}

	return profile{
		In:      ffmpeg.MergeKwArgs(in),
		Out:     ffmpeg.MergeKwArgs(out),
		Filters: MergeFilters(filters),
	}
}

func ReadConfig(name string) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&profiles)
	if err != nil {
		log.Fatal(err)
	}
}

func GetProfile(name string) Cmd {
	var pro profile
	switch name {
	case "default", "stream":
		pro = MergeProfiles("quiet", "stream")
	case "audio":
		pro = MergeProfiles("quiet", "stream", "defaultAudio")
	case "video":
		pro = MergeProfiles("quiet", "stream", "defaultVideo")
	default:
		pro = profiles[name]
	}

	filters := pro.Filters
	if len(filters) == 0 {
		filters = make(Filters)
	}

	cmd := Cmd{
		Filters: filters,
		Output:  NewOutput(pro.Out),
		Input:   NewInput(pro.In),
	}
	return cmd
}
