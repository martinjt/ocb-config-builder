package main

import (
	"embed"

	"github.com/martinjt/ocb-config-builder/pkg/configmapping"
	"gopkg.in/yaml.v3"
)

//go:embed component_mapping.yaml
var componentMappingFile embed.FS

func getComponentMapping() configmapping.ComponentMappingFile {
	componentMappingBytes, _ := componentMappingFile.ReadFile("component_mapping.yaml")
	var componentMapping configmapping.ComponentMappingFile
	yaml.Unmarshal(componentMappingBytes, &componentMapping)
	return componentMapping
}
