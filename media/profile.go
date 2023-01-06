package profile

import (
	"log"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

var profiles = make(map[string]Profile)

type Profile struct {
	Filters map[string]ffmpeg.KwArgs `yaml:"filters"`
	In      ffmpeg.KwArgs            `yaml:"input"`
	Out     ffmpeg.KwArgs            `yaml:"output"`
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

func Get(name string) Profile {
	pro := profiles[name]
	filters := pro.Filters
	if len(filters) == 0 {
		filters = make(map[string]ffmpeg.KwArgs)
	}
	return pro
}
