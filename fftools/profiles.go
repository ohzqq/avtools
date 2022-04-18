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

type ffCfg struct {
	Profiles map[string]CmdArgs
	ProfileList []string
	Padding int
	Output string
	Verbosity string
	ProCfg map[string]*viper.Viper
	Defaults *viper.Viper
}

var Cfg = ffCfg{}

func FFcfg() {
	//cfg := FFCfg{}
	list, profiles := parseProfiles(viper.Get("Profiles"))
	//cfg.proCfg = 
	Cfg.Defaults = viper.Sub("Defaults")
	Cfg.Profiles = profiles
	Cfg.ProfileList = list
	Cfg.Padding = Cfg.Defaults.GetInt("Padding")
	Cfg.Output = Cfg.Defaults.GetString("Output")
	Cfg.Verbosity = Cfg.Defaults.GetString("Verbosity")
	//return &cfg
}

func (c *ffCfg) GetDefaults() *viper.Viper {
	return c.Defaults
}

func parseProfiles(p interface{}) ([]string, map[string]CmdArgs) {
	var name string
	var list []string
	prof := make(map[string]string)
	profiles := make(map[string]CmdArgs)
	proCfg := make(map[string]*viper.Viper)
	for _, pro := range p.([]interface{}) {
		for k, v := range pro.(map[interface{}]interface{}) {
			prof[k.(string)] = v.(string)
			if k.(string) == "Name" {
				name = v.(string)
				list = append(list, v.(string))
			}
		}
		proCfg[name] = viper.Sub(name)
		profiles[name] = CmdArgs(prof)
	}
	return list, profiles
}
