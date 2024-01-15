package main

import (
	"os"

	"github.com/martinjt/ocb-config-builder/pkg/collectorconfig"
	"github.com/martinjt/ocb-config-builder/pkg/configmapping"
	"gopkg.in/yaml.v3"
)

type ocbConfig struct {
	Dist       Dist         `yaml:"dist"`
	Extensions []ModulePath `yaml:"extensions,omitempty"`
	Receivers  []ModulePath `yaml:"receivers,omitempty"`
	Processors []ModulePath `yaml:"processors,omitempty"`
	Exporters  []ModulePath `yaml:"exporters,omitempty"`
	Connectors []ModulePath `yaml:"connectors,omitempty"`
}

type Dist struct {
	Name             string `yaml:"name"`
	Description      string `yaml:"description"`
	Version          string `yaml:"version"`
	CollectorVersion string `yaml:"otelcol_version"`
}

type ModulePath struct {
	GoMod string `yaml:"gomod"`
}

func (config *ocbConfig) addComponent(components collectorconfig.RequiredComponents, componentMapping configmapping.ComponentMappingFile) bool {
	allFound := true
	for _, v := range components.Receivers {
		modPath, allFound := componentMapping.GetConfigType("receiver", v)
		config.Receivers = append(config.Receivers, ModulePath{GoMod: modPath})
		if !allFound {
			return false
		}
	}
	for _, v := range components.Processors {
		modPath, allFound := componentMapping.GetConfigType("processor", v)
		config.Processors = append(config.Processors, ModulePath{GoMod: modPath})
		if !allFound {
			return false
		}
	}
	for _, v := range components.Exporters {
		modPath, allFound := componentMapping.GetConfigType("exporter", v)
		config.Exporters = append(config.Exporters, ModulePath{GoMod: modPath})
		if !allFound {
			return false
		}
	}
	for _, v := range components.Extensions {
		modPath, allFound := componentMapping.GetConfigType("extensions", v)
		config.Extensions = append(config.Extensions, ModulePath{GoMod: modPath})
		if !allFound {
			return false
		}
	}
	for _, v := range components.Connectors {
		modPath, allFound := componentMapping.GetConfigType("connector", v)
		config.Connectors = append(config.Connectors, ModulePath{GoMod: modPath})
		if !allFound {
			return false
		}
	}
	return allFound
}

func (config *ocbConfig) writeConfigToFile(filename string) error {
	configBytes, _ := yaml.Marshal(&config)
	return os.WriteFile(filename, configBytes, 0644)
}
