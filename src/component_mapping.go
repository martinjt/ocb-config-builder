package main

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

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

func getComponentMapping() ComponentMappingFile {
	componentMappingBytes, _ := componentMappingFile.ReadFile("component_mapping.yaml")
	var componentMapping ComponentMappingFile
	yaml.Unmarshal(componentMappingBytes, &componentMapping)
	return componentMapping
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
	return fmt.Sprintf("%v v%v", component.GithubUrl, component.Version)
}
