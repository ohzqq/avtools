package avtools

import (
	"path/filepath"
	"os"
	"fmt"
	"log"
	//"strings"

	//"golang.org/x/exp/slices"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg = ffCfg{ profiles: make(map[string]*Args) }
)

type ffCfg struct {
	pros *viper.Viper
	defaults *viper.Viper
	proList []string
	profiles map[string]*Args
}

// initConfig reads in config file 
func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".avtools" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config/avtools"))
		viper.AddConfigPath(filepath.Join(home, "Sync/code/avtools/tmp/"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yml")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		cfg.pros = viper.Sub("Profiles")

		cfg.defaults = viper.Sub("defaults")


		err := cfg.pros.Unmarshal(&cfg.profiles)
		if err != nil {
			fmt.Printf("unable to decode into struct, %v", err)
		}
	}

	cfg.profiles["default"] = &Args{
		Flags: Flags{Output: "tmp"},
		Padding: "%06d",
		VideoCodec: "copy",
		AudioCodec: "copy",
	}

	for name, pro := range cfg.profiles {
		cfg.proList = append(cfg.proList, name)

		var profile string
		if cfg.defaults.IsSet("profile") {
			profile = cfg.defaults.GetString("profile")
		}

		if name == profile {
			cfg.profiles["default"] = pro
		}
	}
}

func Cfg() ffCfg {
	return cfg
}

func(cfg ffCfg) GetProfile(p string) *Args {
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

func(cfg ffCfg) Profiles() []string {
	return cfg.proList
}

func(cfg ffCfg) DefaultProfile() *Args {
	return cfg.GetProfile("default")
}

