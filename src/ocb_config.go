package main

import (
	"os"

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

func (config *ocbConfig) addComponent(components requiredComponents, componentMapping ComponentMappingFile) {
	for _, v := range components.Receivers {
		config.Receivers = append(config.Receivers, ModulePath{GoMod: componentMapping.GetConfigType("receiver", v)})
	}
	for _, v := range components.Processors {
		config.Processors = append(config.Processors, ModulePath{GoMod: componentMapping.GetConfigType("processor", v)})
	}
	for _, v := range components.Exporters {
		config.Exporters = append(config.Exporters, ModulePath{GoMod: componentMapping.GetConfigType("exporter", v)})
	}
	for _, v := range components.Extensions {
		config.Extensions = append(config.Extensions, ModulePath{GoMod: componentMapping.GetConfigType("extensions", v)})
	}
	for _, v := range components.Connectors {
		config.Connectors = append(config.Connectors, ModulePath{GoMod: componentMapping.GetConfigType("connector", v)})
	}
}

func (config *ocbConfig) writeConfigToFile(filename string) error {
	configBytes, _ := yaml.Marshal(&config)
	return os.WriteFile(filename, configBytes, 0644)
}
