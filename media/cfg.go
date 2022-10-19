package media

import (
	"log"

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
