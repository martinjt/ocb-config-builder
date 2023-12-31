package main

import (
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type collectorConfigFile struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Extensions map[string]interface{} `yaml:"extensions"`
	Connectors map[string]interface{} `yaml:"connectors"`
}

type requiredComponents struct {
	Receivers  []string
	Processors []string
	Exporters  []string
	Extensions []string
	Connectors []string
}

func getRequiredComponentsFromCollectorConfig(filename string) (requiredComponents, error) {

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return requiredComponents{}, err
	}

	var data collectorConfigFile

	// Unmarshal the YAML string into the data map
	yaml.Unmarshal(bytes, &data)

	requiredComponents := requiredComponents{}

	for k := range data.Receivers {
		requiredComponents.Receivers = append(requiredComponents.Receivers, getType(k))
	}
	for k := range data.Processors {
		requiredComponents.Processors = append(requiredComponents.Processors, getType(k))
	}
	for k := range data.Exporters {
		requiredComponents.Exporters = append(requiredComponents.Exporters, getType(k))
	}
	for k := range data.Extensions {
		requiredComponents.Extensions = append(requiredComponents.Extensions, getType(k))
	}
	for k := range data.Connectors {
		requiredComponents.Connectors = append(requiredComponents.Connectors, getType(k))
	}
	requiredComponents.Receivers = slices.Compact(requiredComponents.Receivers)
	requiredComponents.Processors = slices.Compact(requiredComponents.Processors)
	requiredComponents.Exporters = slices.Compact(requiredComponents.Exporters)
	requiredComponents.Extensions = slices.Compact(requiredComponents.Extensions)
	requiredComponents.Connectors = slices.Compact(requiredComponents.Connectors)
	return requiredComponents, nil
}

func getType(fullname string) string {
	return strings.Split(fullname, "/")[0]
}
