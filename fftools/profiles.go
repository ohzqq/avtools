package fftools

import (
	"path/filepath"
	"os"
	"fmt"
	"strings"

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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

//type profile map[string]string
type profile []string
type defaults map[string]interface{}

type FFCfg struct {
	Profiles []string
	Padding int
	Output string
	proCfg *viper.Viper
	defCfg *viper.Viper
}

func Cfg() *FFCfg {
	cfg := FFCfg{}
	cfg.proCfg = viper.Sub("profiles")
	cfg.defCfg = viper.Sub("defaults")
	cfg.Profiles = cfg.listProfiles()
	cfg.Padding = cfg.defCfg.GetInt("padding")
	cfg.Output = cfg.defCfg.GetString("output")
	return &cfg
}

func (ff *FFCfg) listProfiles() []string {
	var profiles []string
	for pro, _ := range viper.GetStringMap("profiles") {
		profiles = append(profiles, pro)
	}
	return profiles
}

func (ff *FFCfg) Profile(pro string) (profile) {
	return profile(ff.proCfg.GetStringSlice(pro))
}
