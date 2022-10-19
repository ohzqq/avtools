package media

import (
	"log"

	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/spf13/viper"
)

var cfg Config

type Config struct {
	Defaults Defaults           `mapstructure:"default"`
	Profiles map[string]Profile `toml:"profile",mapstructure:"profiles"`
}

type Defaults struct {
	Profile   string
	Padding   string
	Output    string
	LogLevel  string
	Overwrite bool
}

type Profile struct {
	PreInput     []string
	PostInput    []string
	VideoCodec   string
	VideoParams  map[string]string
	VideoFilters map[string]string
	AudioCodec   string
	AudioParams  map[string]string
	AudioFilters map[string]string
	Filters      map[string]string
	Ext          string
}

func InitConfig(v *viper.Viper) {
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func Cfg() Config {
	return cfg
}

func (p Profile) FFmpegCmd() *ffmpeg.Cmd {
	ff := ffmpeg.New()

	if len(p.PreInput) > 0 {
		ff.PreInput = p.PreInput
	}

	if len(p.PostInput) > 0 {
		ff.PostInput = p.PostInput
	}

	if p.VideoCodec != "" {
		ff.SetVideoCodec(p.VideoCodec)
	}

	if len(p.VideoParams) > 0 {
		for k, v := range p.VideoParams {
			ff.AppendVideoParam(k, v)
		}
	}

	if len(p.VideoFilters) > 0 {
		for k, v := range p.VideoFilters {
			f := k
			if k == "misc" {
				f = ""
			}

			filter := ffmpeg.NewFilter(f)
			filter.Set(v)
			ff.AppendVideoFilter(filter)
		}
	}

	if p.AudioCodec != "" {
		ff.SetAudioCodec(p.AudioCodec)
	}

	if len(p.AudioParams) > 0 {
		for k, v := range p.AudioParams {
			ff.AppendAudioParam(k, v)
		}
	}

	if len(p.AudioFilters) > 0 {
		for k, v := range p.AudioFilters {
			f := k
			if k == "misc" {
				f = ""
			}

			filter := ffmpeg.NewFilter(f)
			filter.Set(v)
			ff.AppendAudioFilter(filter)
		}
	}

	if len(p.Filters) > 0 {
		for k, v := range p.Filters {
			f := k
			if k == "misc" {
				f = ""
			}

			filter := ffmpeg.NewFilter(f)
			filter.Set(v)
			ff.AppendFilter(filter)
		}
	}

	if p.Ext != "" {
	}

	return ff
}

func (d Defaults) HasPadding() bool {
	return d.Padding != ""
}

func (d Defaults) HasLogLevel() bool {
	return d.LogLevel != ""
}

func (d Defaults) HasOutput() bool {
	return d.Output != ""
}

func (d Defaults) HasProfile() bool {
	return d.Profile != ""
}
