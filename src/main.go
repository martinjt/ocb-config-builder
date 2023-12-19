package main

import (
	"embed"
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
)

type CollectorConfigFile struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Extensions map[string]interface{} `yaml:"extensions"`
	Connectors map[string]interface{} `yaml:"connectors"`
}

type Config struct {
	Dist       Dist         `yaml:"dist"`
	Extensions []ModulePath `yaml:"extensions"`
	Receivers  []ModulePath `yaml:"receivers"`
	Processors []ModulePath `yaml:"processors"`
	Exporters  []ModulePath `yaml:"exporters"`
	Connectors []ModulePath `yaml:"connectors"`
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

type ComponentMapping struct {
	GithubUrl string `yaml:"github_url"`
	Version   string `yaml:"version"`
}

type ComponentMappingFile struct {
	Receivers  map[string]ComponentMapping `yaml:"receivers"`
	Processors map[string]ComponentMapping `yaml:"processors"`
	Exporters  map[string]ComponentMapping `yaml:"exporters"`
	Extensions map[string]ComponentMapping `yaml:"extensions"`
	Connectors map[string]ComponentMapping `yaml:"connectors"`
}

//go:embed component_mapping.yaml
var componentMappingFile embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide an input file")
		return
	}
	inputFile := os.Args[1]

	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	componentMappingBytes, _ := componentMappingFile.ReadFile("component_mapping.yaml")

	var componentMapping ComponentMappingFile
	yaml.Unmarshal(componentMappingBytes, &componentMapping)

	var data CollectorConfigFile

	// Unmarshal the YAML string into the data map
	yaml.Unmarshal(bytes, &data)

	var receivers []string
	var processors []string
	var exporters []string
	var extensions []string
	var connectors []string

	for k := range data.Receivers {
		receivers = append(receivers, getType(k))
	}
	for k := range data.Processors {
		processors = append(processors, getType(k))
	}
	for k := range data.Exporters {
		exporters = append(exporters, getType(k))
	}
	for k := range data.Extensions {
		extensions = append(extensions, getType(k))
	}
	for k := range data.Connectors {
		connectors = append(connectors, getType(k))
	}

	config := Config{}
	config.Dist.Name = "otelcol-custom"
	config.Dist.Description = "Custom distribution of the OpenTelemetry Collector"
	config.Dist.Version = "0.90.1"
	config.Dist.CollectorVersion = "0.90.1"

	for _, v := range slices.Compact(extensions) {
		config.Extensions = append(config.Extensions, ModulePath{GoMod: componentMapping.GetConfigType("extensions", v)})
	}
	for _, v := range slices.Compact(receivers) {
		config.Receivers = append(config.Receivers, ModulePath{GoMod: componentMapping.GetConfigType("receiver", v)})
	}
	for _, v := range slices.Compact(processors) {
		config.Processors = append(config.Processors, ModulePath{GoMod: componentMapping.GetConfigType("processor", v)})
	}
	for _, v := range slices.Compact(exporters) {
		config.Exporters = append(config.Exporters, ModulePath{GoMod: componentMapping.GetConfigType("exporter", v)})
	}
	for _, v := range slices.Compact(connectors) {
		config.Connectors = append(config.Connectors, ModulePath{GoMod: componentMapping.GetConfigType("connector", v)})
	}

	configBytes, _ := yaml.Marshal(&config)

	os.WriteFile("dist.yaml", configBytes, 0644)
}

func getType(fullname string) string {
	return strings.Split(fullname, "/")[0]
}

func (c *ComponentMappingFile) GetConfigType(componentType string, componentTypeName string) string {
	var component ComponentMapping
	var found bool = false
	switch componentType {
	case "receiver":
		component, found = c.Receivers[componentTypeName]
	case "processor":
		component, found = c.Processors[componentTypeName]
	case "exporter":
		component, found = c.Exporters[componentTypeName]
	case "extensions":
		component, found = c.Extensions[componentTypeName]
	case "connector":
		component, found = c.Connectors[componentTypeName]
	}
	if !found {
		fmt.Printf("Component not found: %v:%v \n", componentType, componentTypeName)
		return ""
	}
	return fmt.Sprintf("%s  v%s", component.GithubUrl, component.Version)
}
