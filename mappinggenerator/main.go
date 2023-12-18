package main

import (
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
	GithubUrl        string `yaml:"github_url"`
	ConfigTypeString string `yaml:"config_type"`
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

	var receivers []ComponentMapping
	var processors []ComponentMapping
	var exporters []ComponentMapping
	var extensions []ComponentMapping
	var connectors []ComponentMapping

	tree.Files().ForEach(func(f *object.File) error {
		if filepath.Base(f.Name) == "metadata.yaml" {
			splitpath := strings.Split(f.Name, "/")
			componentType := splitpath[0]
			componentDir := splitpath[1]

			metadata := Metadata{}

			metadataContents, _ := f.Contents()

			yaml.Unmarshal([]byte(metadataContents), &metadata)

			componentMapping := ComponentMapping{
				GithubUrl:        "github.com/open-telemetry/opentelemetry-collector-contrib/" + componentType + "/" + componentDir,
				ConfigTypeString: metadata.ConfigType,
			}

			switch componentType {
			case "receivers":
				receivers = append(receivers, componentMapping)
			case "processors":
				processors = append(processors, componentMapping)
			case "exporters":
				exporters = append(exporters, componentMapping)
			case "extensions":
				extensions = append(extensions, componentMapping)
			case "connectors":
				connectors = append(connectors, componentMapping)
			default:
				panic("Unknown component type: " + componentType)
			}
		}
		return nil
	})
}
