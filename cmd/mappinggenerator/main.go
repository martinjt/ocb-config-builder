package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/martinjt/ocb-config-builder/pkg/configmapping"
	"github.com/martinjt/ocb-config-builder/pkg/mappinggenerator"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	ConfigType string `yaml:"type"`
}

const (
	CONTRIB_REPO      = "github.com/open-telemetry/opentelemetry-collector-contrib"
	COLLECTOR_VERSION = "0.95.0"
)

var opts struct {
	ContribVersion string `short:"v" long:"version" description:"The version of the contrib repo to generate a mapping file for"`
	Verbose        bool   `short:"d" long:"debug" description:"Enable debug logging"`
}

func main() {

	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	if opts.ContribVersion == "" {
		opts.ContribVersion = COLLECTOR_VERSION
	}

	config := configmapping.ComponentMappingFile{
		Receivers:  getCoreReceiverMapping(COLLECTOR_VERSION),
		Processors: getCoreProcessorMapping(COLLECTOR_VERSION),
		Exporters:  getCoreExporterMapping(COLLECTOR_VERSION),
		Extensions: getCoreExtensionMapping(COLLECTOR_VERSION),
		Connectors: getCoreConnectorMapping(COLLECTOR_VERSION),
	}

	contribRepoMappings := mappinggenerator.GenerateMappingFileForRepo(CONTRIB_REPO, opts.ContribVersion, opts.Verbose)

	config.MergeMappingFiles(contribRepoMappings)

	configBytes, _ := yaml.Marshal(&config)

	os.WriteFile("cmd/configgenerator/component_mapping.yaml", configBytes, 0644)
}
