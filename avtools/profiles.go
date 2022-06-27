package avtools

import (
	"path/filepath"
	"os"
	"fmt"
	"log"
	//"strings"

	"golang.org/x/exp/slices"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg = ffCfg{
		profiles: make(pros),
		defaults: defaults{
			Output: "tmp",
			Padding: "%06d",
		},
	}
)

type ffCfg struct {
	proList []string
	profiles pros
	defaults
}

func Cfg() ffCfg {
	return cfg
}

func(cfg ffCfg) Profiles() []string {
	return cfg.proList
}

func(cfg ffCfg) DefaultProfile() string {
	def := "base"
	if d := cfg.defaults.Profile; slices.Contains(cfg.Profiles(), d) {
		def = d
	}
	return def
}

func(cfg ffCfg) GetProfile(p string) *Args {
	return cfg.profiles[p]
}

type defaults struct {
	Output string
	LogLevel string
	Overwrite bool
	Profile string
	Padding string
}

type pros map[string]*Args

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
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

	viper.AutomaticEnv() // read in environment variables that match

	var (
		padding = "%06d"
		output = "tmp"
		overwrite bool
		profile string
		log string
	)
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		defaults := viper.Sub("defaults")
		pros := viper.Sub("Profiles")
		if defaults.IsSet("padding") {
			padding = defaults.GetString("padding")
		}
		if defaults.IsSet("output") {
			output = defaults.GetString("output")
		}
		if defaults.IsSet("profile") {
			profile = defaults.GetString("profile")
		}
		if defaults.IsSet("loglevel") {
			log = defaults.GetString("loglevel")
		}
		if defaults.IsSet("overwrite") {
			overwrite = defaults.GetBool("overwrite")
		}
		err := pros.Unmarshal(&cfg.profiles)
		if err != nil {
			fmt.Printf("unable to decode into struct, %v", err)
		}
	}

	cfg.profiles["base"] = &Args{
		VideoCodec: "copy",
		AudioCodec: "copy",
	}

	for name, pro := range cfg.profiles {
		cfg.proList = append(cfg.proList, name)

		if pro.Output == "" {
			pro.Output = output
		}

		if pro.Padding == "" {
			pro.Padding = padding
		}

		if pro.LogLevel == "" {
			pro.LogLevel = log
		}

		if !pro.Overwrite {
			pro.Overwrite = overwrite
		}

		if name == profile {
			cfg.profiles["default"] = pro
		}
	}
}

func (p pros) List() []string {
	var list []string
	for pro, _ := range p {
		list = append(list, pro)
	}
	return list
}

