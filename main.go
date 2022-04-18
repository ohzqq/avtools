package main

import (
	//"fmt"
	"log"

	"github.com/ohzqq/fftools/cmd"
	"github.com/ohzqq/fftools/fftools"
)

func main() {
	fftools.InitConfig()
	fftools.FFcfg()
	cli := cmd.NewCli()
	//fmt.Printf("%v", fftools.Cfg.Profiles)
	err := cli.Cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//cfg := fftools.Cfg()
	//keys := fftools.Profiles()
	//defaults := fftools.Defaults()
	//fmt.Printf("%v", fftools.CfgFile)
	//pro := cfg.Profile("convert")
	//fmt.Printf("%v", pro.String())
	//fmt.Printf("%v", defaults)
	//fmt.Printf("%s-%02d.mkv", "output", 0)
	//cmd.Execute()
}

