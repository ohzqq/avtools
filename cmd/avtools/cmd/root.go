package cmd

import (
	"fmt"
	"mime"
	"os"

	"github.com/ohzqq/avtools/media"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type fmtBoolFlags struct {
	Meta  bool
	Cue   bool
	Cover bool
}

func (flags fmtBoolFlags) DumpMeta(m *media.Media) []media.Cmd {
	var cmds []media.Cmd

	if flags.Cue {
		c := m.SaveMetaFmt("cue")
		cmds = append(cmds, c)
	}

	if flags.Meta {
		c := m.SaveMetaFmt("ffmeta")
		cmds = append(cmds, c)
	}

	return cmds
}

func (flags fmtBoolFlags) Extract(m *media.Media) []media.Cmd {
	var cmds []media.Cmd

	if flags.Cue {
		c := m.SaveMetaFmt("cue")
		cmds = append(cmds, c)
	}

	if flags.Meta {
		c := m.SaveMetaFmt("ffmeta")
		cmds = append(cmds, c)
	}

	if flags.Cover {
		ff := m.ExtractCover()
		cmds = append(cmds, ff)
	}

	return cmds
}

type fmtStringFlags struct {
	Meta  string
	Cue   string
	Cover string
}

func (flags fmtStringFlags) LoadMeta(m *media.Media) {
	switch {
	case flags.Meta != "":
		m.LoadIni(flags.Meta)
	case flags.Cue != "":
		m.LoadCue(flags.Cue)
	}
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	mime.AddExtensionType(".ini", "text/plain")
	mime.AddExtensionType(".cue", "text/plain")
	mime.AddExtensionType(".m4b", "audio/mp4")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cmd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
