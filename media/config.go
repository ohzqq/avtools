//go:build exclude

package media

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	cfg = AVcfg{profiles: make(map[string]*Args)}
)

type AVcfg struct {
	pros     *viper.Viper
	defaults defaults
	ProList  []string
	profiles map[string]*Args
}

type defaults struct {
	Profile   string
	Padding   string
	Output    string
	LogLevel  string
	Overwrite bool
}

func Cfg() AVcfg {
	return cfg
}

func InitCfg() {
	cfg.defaults = defaults{
		Padding:   "%06d",
		Profile:   "default",
		Output:    "tmp",
		LogLevel:  "error",
		Overwrite: false,
	}

	cfg.ProList = append(cfg.ProList, "default")
	cfg.profiles["default"] = &Args{
		VideoCodec: "copy",
		AudioCodec: "copy",
	}
}

func CfgProfiles(def, profiles *viper.Viper) {
	var err error
	err = def.Unmarshal(&cfg.defaults)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	cfg.pros = profiles
	err = cfg.pros.Unmarshal(&cfg.profiles)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	for name, pro := range cfg.profiles {
		cfg.ProList = append(cfg.ProList, name)

		var profile string
		if def.IsSet("profile") {
			profile = def.GetString("profile")
		}

		if name == profile {
			cfg.profiles["default"] = pro
		}
	}
}

func (cfg AVcfg) GetProfile(p string) *Args {
	pro := cfg.profiles[p]

	if pro.Padding == "" {
		pro.Padding = cfg.defaults.Padding
	}

	if pro.Output == "" {
		pro.Output = cfg.defaults.Output
	}

	if pro.LogLevel == "" {
		pro.LogLevel = cfg.defaults.LogLevel
	}

	if !pro.Overwrite {
		pro.Overwrite = cfg.defaults.Overwrite
	}

	return pro
}

func (cfg AVcfg) GetDefault(p string) string {
	switch p {
	case "padding":
		return cfg.defaults.Padding
	case "output":
		return cfg.defaults.Output
	case "profile":
		return cfg.defaults.Profile
	case "logLevel":
		return cfg.defaults.LogLevel
	}
	return ""
}

func (cfg AVcfg) OverwriteDefault() bool {
	return cfg.defaults.Overwrite
}

func (cfg AVcfg) Profiles() []string {
	return cfg.ProList
}

func (cfg AVcfg) DefaultProfile() *Args {
	return cfg.GetProfile("default")
}
