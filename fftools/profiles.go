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

//type profile map[string]string
type defaults map[string]interface{}

type FFCfg struct {
	Profiles map[string]profile
	ProfileList []string
	Padding int
	Output string
	proCfg []interface{}
	defCfg *viper.Viper
}

type Profile struct {
	Pre string
	VideoCodec string
	VideoFilters string
	VideoParams string
	AudioCodec string
	AudioFilters string
	AudioParams string
	FilterComplex string
	Post string
}

func Cfg() *FFCfg {
	cfg := FFCfg{}
	list, profiles := parseProfiles(viper.Get("profiles"))
	//cfg.proCfg = 
	cfg.defCfg = viper.Sub("defaults")
	cfg.Profiles = profiles
	cfg.ProfileList = list
	cfg.Padding = cfg.defCfg.GetInt("padding")
	cfg.Output = cfg.defCfg.GetString("output")
	return &cfg
}

type profile map[string]string

func parseProfiles(p interface{}) ([]string, map[string]profile) {
	var name string
	var list []string
	prof := make(map[string]string)
	profiles := make(map[string]profile)
	for _, pro := range p.([]interface{}) {
		for k, v := range pro.(map[interface{}]interface{}) {
			prof[k.(string)] = v.(string)
			if k.(string) == "Name" {
				name = v.(string)
				list = append(list, v.(string))
			}
		}
		profiles[name] = profile(prof)
	}
	fmt.Printf("%v", profiles)
	return list, profiles
}

//func (ff *FFCfg) listProfiles() []string {
//  var profiles []string
//  for pro, _ := range viper.GetStringMap("profiles") {
//    profiles = append(profiles, pro)
//  }
//  return profiles
//}

//func (ff *FFCfg) Profile(pro string) (profile) {
	//return profile(ff.proCfg.GetStringSlice(pro))
//}
