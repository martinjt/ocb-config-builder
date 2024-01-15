package configmapping

import "fmt"

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
