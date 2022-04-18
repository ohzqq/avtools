package fftools

import (
	"path/filepath"
	"os"
	//"fmt"
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

type Profile map[string]string

type FFCfg struct {
	Profiles map[string]Profile
	ProfileList []string
	Padding int
	Output string
}

func Cfg() *FFCfg {
	cfg := FFCfg{}
	list, profiles := parseProfiles(viper.Get("Profiles"))
	//cfg.proCfg = 
	defCfg := viper.Sub("Defaults")
	cfg.Profiles = profiles
	cfg.ProfileList = list
	cfg.Padding = defCfg.GetInt("Padding")
	cfg.Output = defCfg.GetString("Output")
	return &cfg
}

func parseProfiles(p interface{}) ([]string, map[string]Profile) {
	var name string
	var list []string
	prof := make(map[string]string)
	profiles := make(map[string]Profile)
	for _, pro := range p.([]interface{}) {
		for k, v := range pro.(map[interface{}]interface{}) {
			prof[k.(string)] = v.(string)
			if k.(string) == "Name" {
				name = v.(string)
				list = append(list, v.(string))
			}
		}
		profiles[name] = Profile(prof)
	}
	return list, profiles
}
