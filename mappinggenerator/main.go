package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	ConfigType string `yaml:"type"`
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

func main() {

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           "https://github.com/open-telemetry/opentelemetry-collector-contrib.git",
		NoCheckout:    true,
		ReferenceName: "main",
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	headCommit, err := r.Head()
	if err != nil {
		panic(err)
	}

	commit, err := r.CommitObject(headCommit.Hash())
	if err != nil {
		panic(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}

	receivers := getCoreReceiverMapping()
	processors := getCoreProcessorMapping()
	exporters := getCoreExporterMapping()
	extensions := make(map[string]ComponentMapping)
	connectors := make(map[string]ComponentMapping)

	tree.Files().ForEach(func(f *object.File) error {
		if filepath.Base(f.Name) == "metadata.yaml" {
			splitpath := strings.Split(f.Name, "/")
			componentType := splitpath[0]
			componentDir := splitpath[1]

			metadata := Metadata{}

			metadataContents, _ := f.Contents()

			yaml.Unmarshal([]byte(metadataContents), &metadata)

			componentMapping := ComponentMapping{
				GithubUrl: "github.com/open-telemetry/opentelemetry-collector-contrib/" + componentType + "/" + componentDir,
				Version:   "0.90.1",
			}

			switch componentType {
			case "receiver":
				receivers[metadata.ConfigType] = componentMapping
			case "processor":
				processors[metadata.ConfigType] = componentMapping
			case "exporter":
				exporters[metadata.ConfigType] = componentMapping
			case "extension":
				extensions[metadata.ConfigType] = componentMapping
			case "connector":
				connectors[metadata.ConfigType] = componentMapping
			default:
				fmt.Println("Unknown component type: " + componentType)
			}
		}
		return nil
	})

	config := ComponentMappingFile{
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Extensions: extensions,
		Connectors: connectors,
	}

	configBytes, _ := yaml.Marshal(&config)

	os.WriteFile("../src/component_mapping.yaml", configBytes, 0644)
}

func getCoreProcessorMapping() map[string]ComponentMapping {
	return map[string]ComponentMapping{
		"batch": {
			GithubUrl: "go.opentelemetry.io/collector/processor/batchprocessor",
			Version:   "0.90.1",
		},
	}
}

func getCoreExporterMapping() map[string]ComponentMapping {
	return map[string]ComponentMapping{
		"logging": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/loggingexporter",
			Version:   "0.90.1",
		},
		"otlp": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/otlpexporter",
			Version:   "0.90.1",
		},
		"otlphttp": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/otlphttpexporter",
			Version:   "0.90.1",
		},
	}
}

func getCoreReceiverMapping() map[string]ComponentMapping {
	return map[string]ComponentMapping{
		"otlp": {
			GithubUrl: "go.opentelemetry.io/collector/receiver/otlpreceiver",
			Version:   "0.90.1",
		},
	}
}
