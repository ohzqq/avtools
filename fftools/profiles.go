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

type ffCfg struct {
	Profiles pros
	//Profiles map[string]CmdArgs
	ProfileList []string
	//ProCfg *viper.Viper
	Defaults defaults
}

type defaults struct {
	Padding bool
	Output string
	Verbosity string
	Overwrite bool
}

var Cfg = ffCfg{}

func FFcfg() {
	viper.Sub("Defaults").Unmarshal(&Cfg.Defaults)
	proCfg := viper.Sub("Profiles")
	Cfg.Profiles = parseProfiles(proCfg)
	Cfg.ProfileList = Cfg.profileList()
}

func (c *ffCfg) profileList() []string {
	var list []string
	for pro, _ := range c.Profiles {
		list = append(list, pro)
	}
	return list
}

//func (c *ffCfg) GetDefaults() *viper.Viper {
//  return c.Defaults
//}

type Profile struct {
	Pre string
	Input string
	Post string
	VideoCodec string
	VideoParams string
	VideoFilters string
	AudioCodec string
	AudioParams string
	AudioFilters string
	FilterCompex string
	Verbosity string
	Output string
	Padding bool
	Ext string
	Overwrite bool
}

type pros map[string]Profile

func parseProfiles(v *viper.Viper) pros {
	profiles := make(pros)
	err := v.Unmarshal(&profiles)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return profiles
}
