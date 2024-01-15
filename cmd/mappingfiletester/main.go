package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/martinjt/ocb-config-builder/pkg/mappinggenerator"
	"gopkg.in/yaml.v3"
)

var opts struct {
	Repo    string `short:"r" long:"repo" description:"The repo to generate a mapping file for"`
	Verbose bool   `short:"d" long:"debug" description:"Enable debug logging"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	contribRepoMappings := mappinggenerator.GenerateMappingFileForRepo(opts.Repo, "", opts.Verbose)

	configBytes, _ := yaml.Marshal(&contribRepoMappings)

	println(string(configBytes))
}
