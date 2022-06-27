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
	profiles pros
	defaults
}

func Cfg() ffCfg {
	return cfg
}

func(cfg ffCfg) Profiles() []string {
	var list []string
	for pro, _ := range cfg.profiles {
		list = append(list, pro)
	}
	return list
}

func(cfg ffCfg) DefaultProfile() string {
	def := "base"
	if d := cfg.defaults.Profile; slices.Contains(cfg.Profiles(), d) {
		def = d
	}
	return def
}

func(cfg ffCfg) GetProfile(p string) *CmdArgs {
	return cfg.profiles[p]
}

type defaults struct {
	Output string
	LogLevel string
	Overwrite bool
	Profile string
	Padding string
}

type pros map[string]*CmdArgs

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

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		viper.Sub("Defaults").Unmarshal(&cfg.defaults)
		err := viper.Sub("Profiles").Unmarshal(&cfg.profiles)
		fmt.Printf("%+v\n", cfg.Pros["gif"])
		if err != nil {
			fmt.Printf("unable to decode into struct, %v", err)
		}
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	Cfg().Pros["base"] = Arg{
		VideoCodec: "copy",
		AudioCodec: "copy",
	}
}

func (p pros) List() []string {
	var list []string
	for pro, _ := range p {
		list = append(list, pro)
	}
	return list
}

