package avtools

import (
	//"path/filepath"
	//"os"
	"fmt"
	//"log"
	//"strings"

	//"golang.org/x/exp/slices"
	"github.com/spf13/viper"
)

var (
	cfg = AVcfg{ profiles: make(map[string]*Args) }
)

type AVcfg struct {
	pros *viper.Viper
	defaults *viper.Viper
	ProList []string
	profiles map[string]*Args
}

func InitProfiles(defaults, profiles *viper.Viper) {
	cfg.defaults = defaults
	cfg.pros = profiles

	err := cfg.pros.Unmarshal(&cfg.profiles)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	cfg.profiles["default"] = &Args{
		Flags: Flags{Output: "tmp"},
		Padding: "%06d",
		VideoCodec: "copy",
		AudioCodec: "copy",
	}

	for name, pro := range cfg.profiles {
		cfg.ProList = append(cfg.ProList, name)

		var profile string
		if cfg.defaults.IsSet("profile") {
			profile = cfg.defaults.GetString("profile")
		}

		if name == profile {
			cfg.profiles["default"] = pro
		}
	}
}

func Cfg() AVcfg {
	return cfg
}

func(cfg AVcfg) GetProfile(p string) *Args {
	pro := cfg.profiles[p]

	if pro.Padding == "" {
		pro.Padding = cfg.defaults.GetString("padding")
	}
	if pro.Output == "" {
		pro.Output = cfg.defaults.GetString("output")
	}
	if pro.LogLevel == "" {
		pro.LogLevel = cfg.defaults.GetString("loglevel")
	}
	if !pro.Overwrite {
		pro.Overwrite = cfg.defaults.GetBool("overwrite")
	}
	return pro
}

func(cfg AVcfg) GetDefault(p string) string {
	if p != "overwrite" {
		return cfg.defaults.GetString(p)
	}
	return ""
}

func(cfg AVcfg) OverwriteDefault() bool {
	return cfg.defaults.GetBool("overwrite")
}

func(cfg AVcfg) Profiles() []string {
	return cfg.ProList
}

func(cfg AVcfg) DefaultProfile() *Args {
	return cfg.GetProfile("default")
}

