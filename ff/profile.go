package ff

import (
	"log"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

var profiles = make(map[string]profile)

type profile struct {
	Filters Filters       `yaml:"filters"`
	In      ffmpeg.KwArgs `yaml:"input"`
	Out     ffmpeg.KwArgs `yaml:"output"`
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
	pro := profiles[name]
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
