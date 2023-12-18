package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Dist       Dist         `yaml:"dist"`
	Extensions []ModulePath `yaml:"extensions"`
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

	var data map[string]interface{}

	// Unmarshal the YAML string into the data map
	yaml.Unmarshal(bytes, &data)

	var receivers []string
	var processors []string
	var exporters []string
	var extensions []string
	var connectors []string

	for k, _ := range data["receivers"].(map[interface{}]interface{}) {
		receivers = append(receivers, getType(k.(string)))
	}
	for k, _ := range data["processors"].(map[interface{}]interface{}) {
		processors = append(processors, getType(k.(string)))
	}
	for k, _ := range data["exporters"].(map[interface{}]interface{}) {
		exporters = append(exporters, getType(k.(string)))
	}
	for k, _ := range data["extensions"].(map[interface{}]interface{}) {
		extensions = append(extensions, getType(k.(string)))
	}
	for k, _ := range data["connectors"].(map[interface{}]interface{}) {
		connectors = append(connectors, getType(k.(string)))
	}

	fmt.Println("Receivers:")
	for _, v := range slices.Compact(receivers) {
		fmt.Println(v)
	}
	fmt.Println()
	fmt.Println("Processors:")
	for _, v := range slices.Compact(processors) {
		fmt.Println(v)
	}
	fmt.Println()
	fmt.Println("Exporters:")
	for _, v := range slices.Compact(exporters) {
		fmt.Println(v)
	}
	fmt.Println()
	fmt.Println("Extensions:")
	for _, v := range slices.Compact(extensions) {
		fmt.Println(v)
	}
	fmt.Println()
	fmt.Println("Connectors:")
	for _, v := range slices.Compact(connectors) {
		fmt.Println(v)
	}

	config := Config{}
	config.Dist.Name = "otelcol-custom"
	config.Dist.Description = "Custom distribution of the OpenTelemetry Collector"
	config.Dist.Version = "0.90.1"
	config.Dist.CollectorVersion = "0.90.1"

	for _, v := range slices.Compact(extensions) {
		config.Extensions = append(config.Extensions, ModulePath{GoMod: asCore("extensions", v)})
	}

	configBytes, _ := yaml.Marshal(&config)

	os.WriteFile("dist.yaml", configBytes, 0644)
}

func getType(fullname string) string {
	return strings.Split(fullname, "/")[0]
}

func asCore(section string, moduleName string) string {
	return fmt.Sprintf("github.com/open-telemetry/opentelemetry-collector/%v/%v", section, moduleName)
}
