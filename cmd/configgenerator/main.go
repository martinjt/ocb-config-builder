package main

import (
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/martinjt/ocb-config-builder/pkg/collectorconfig"
	"github.com/martinjt/ocb-config-builder/pkg/mappinggenerator"
)

var opts struct {
	ConfigOutput     string `short:"o" long:"output" description:"Output file for the generated config" required:"true" env:"CONFIG_OUTPUT_LOCATION"`
	InputConfig      string `short:"i" long:"input" description:"Input file for the collector config" required:"true" env:"CONFIG_INPUT_LOCATION"`
	CollectorVersion string `short:"v" long:"version" description:"Version of the contrib repo to generate a mapping file for" required:"true" env:"COLLECTOR_VERSION"`
	AdditionalRepos  string `short:"r" long:"additional-repos" description:"Additional repos to generate a mapping file for" env:"ADDITIONAL_REPOS"`
	Verbose          bool   `short:"d" long:"debug" description:"Enable debug logging"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	requiredComponents, err := collectorconfig.GetRequiredComponentsFromCollectorConfig(opts.InputConfig)
	if err != nil {
		panic(err)
	}

	componentMapping := getComponentMapping()

	if opts.AdditionalRepos != "" {
		additionalRepos := strings.Split(opts.AdditionalRepos, ",")
		for _, repo := range additionalRepos {
			additionalRepoMappings := mappinggenerator.GenerateMappingFileForRepo(repo, "", opts.Verbose)
			componentMapping.MergeMappingFiles(additionalRepoMappings)
		}
	}

	config := ocbConfig{}
	config.Dist.Name = "otelcol-custom"
	config.Dist.Description = "Custom distribution of the OpenTelemetry Collector"
	config.Dist.Version = opts.CollectorVersion
	config.Dist.CollectorVersion = opts.CollectorVersion

	if !config.addComponent(requiredComponents, componentMapping) {
		panic(fmt.Errorf("failed to find all components in mapping file"))
	}
	err = config.writeConfigToFile(opts.ConfigOutput)
	if err != nil {
		fmt.Println(err)
		return
	}
}
