package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/martinjt/ocb-config-builder/pkg/collectorconfig"
)

var opts struct {
	ConfigOutput string `short:"o" long:"output" description:"Output file for the generated config" required:"true" env:"CONFIG_OUTPUT_LOCATION"`
	InputConfig  string `short:"i" long:"input" description:"Input file for the collector config" required:"true" env:"CONFIG_INPUT_LOCATION"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	requiredComponents, err := collectorconfig.GetRequiredComponentsFromCollectorConfig(opts.InputConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	componentMapping := getComponentMapping()

	config := ocbConfig{}
	config.Dist.Name = "otelcol-custom"
	config.Dist.Description = "Custom distribution of the OpenTelemetry Collector"
	config.Dist.Version = "0.90.1"
	config.Dist.CollectorVersion = "0.90.1"

	config.addComponent(requiredComponents, componentMapping)
	err = config.writeConfigToFile(opts.ConfigOutput)
	if err != nil {
		fmt.Println(err)
		return
	}
}
