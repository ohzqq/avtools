package fftools

import (
	"path/filepath"
	"os"
	"fmt"
	//"strings"

	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

var CfgFile string

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fftools" (without extension).
		//viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, "Sync/code/fftools/tmp/"))
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fftools.yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

var Cfg = ffCfg{}

type ffCfg struct {
	Profiles pros
	Defaults defaults
}

type defaults struct {
	Padding bool
	Output string
	Verbosity string
	Overwrite bool
	Profile string
}

func FFcfg() {
	viper.Sub("Defaults").Unmarshal(&Cfg.Defaults)
	proCfg := viper.Sub("Profiles")
	Cfg.Profiles = parseProfiles(proCfg)
}

type pros map[string]CmdArgs

func (p pros) List() []string {
	var list []string
	for pro, _ := range p {
		list = append(list, pro)
	}
	return list
}

func parseProfiles(v *viper.Viper) pros {
	profiles := make(pros)
	err := v.Unmarshal(&profiles)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return profiles
}
