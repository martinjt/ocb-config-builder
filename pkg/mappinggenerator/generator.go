package mappinggenerator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/martinjt/ocb-config-builder/pkg/configmapping"
	"gopkg.in/yaml.v3"
)

type metadata struct {
	ConfigType string         `yaml:"type"`
	Status     metadataStatus `yaml:"status"`
	Parent     string         `yaml:"parent"`
}

type metadataStatus struct {
	Class string `yaml:"class"`
}

func GenerateMappingFileForRepo(repoUrl string, version string, verbose bool) configmapping.ComponentMappingFile {
	repoUrl, _ = strings.CutSuffix(repoUrl, "/")
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           fmt.Sprintf("https://%s.git", repoUrl),
		NoCheckout:    true,
		ReferenceName: "main",
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	headCommit, err := repo.Head()
	if err != nil {
		panic(err)
	}

	commit, err := repo.CommitObject(headCommit.Hash())
	if err != nil {
		panic(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}

	receivers := make(map[string]configmapping.ComponentMapping)
	processors := make(map[string]configmapping.ComponentMapping)
	exporters := make(map[string]configmapping.ComponentMapping)
	extensions := make(map[string]configmapping.ComponentMapping)
	connectors := make(map[string]configmapping.ComponentMapping)

	versions := make(map[string]string)
	if version == "" {
		versions = getVersionMap(repo)
	}

	tree.Files().ForEach(func(f *object.File) error {
		if filepath.Base(f.Name) == "metadata.yaml" {
			metadata := metadata{}
			metadataContents, _ := f.Contents()

			yaml.Unmarshal([]byte(metadataContents), &metadata)
			componentType := metadata.Status.Class
			if componentType == "" {
				if verbose {
					fmt.Printf("Component Type for '%s' was empty\n", f.Name)
				}
				return nil
			}

			if metadata.Parent != "" {
				if verbose {
					fmt.Printf("Component Type: %s in '%s' is a child component\n", metadata.ConfigType, f.Name)
				}
				return nil
			}

			if componentType != "receiver" &&
				componentType != "processor" &&
				componentType != "exporter" &&
				componentType != "extension" &&
				componentType != "connector" {
				if verbose {
					fmt.Printf("Component Type: %s in '%s' not external component\n", componentType, f.Name)
				}
				return nil
			}

			directoryPath := filepath.Dir(f.Name)

			componentMapping := configmapping.ComponentMapping{
				GithubUrl: fmt.Sprintf("%s/%s", repoUrl, directoryPath),
			}

			if version == "" {
				componentMapping.Version = versions[directoryPath]
			} else {
				componentMapping.Version = version
			}

			switch metadata.Status.Class {
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
				if verbose {
					fmt.Println("Unknown component type: " + metadata.Status.Class + " Filename: " + f.Name)
				}
			}
		}
		return nil
	})

	config := configmapping.ComponentMappingFile{
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Extensions: extensions,
		Connectors: connectors,
	}

	return config

}

type componentInfo struct {
	Versions []*semver.Version
}

func getVersionMap(repo *git.Repository) map[string]string {
	semverVersions := make(map[string]*componentInfo)
	if tags, err := repo.Tags(); err == nil {
		tags.ForEach(func(t *plumbing.Reference) error {
			if !t.Name().IsTag() {
				return nil
			}
			tagName := strings.ToLower(t.Name().Short())
			tagSplit := strings.Split(tagName, "/")
			if len(tagSplit) < 2 {
				return nil
			}
			tagPath := ""
			for segmentNo := 0; segmentNo < len(tagSplit)-1; segmentNo++ {
				if segmentNo == 0 {
					tagPath = tagSplit[segmentNo]
				} else {
					tagPath = fmt.Sprintf("%s/%s", tagPath, tagSplit[segmentNo])
				}
			}
			componentVersion, _ := semver.NewVersion(tagSplit[(len(tagSplit) - 1)])
			if _, ok := semverVersions[tagPath]; !ok {
				semverVersions[tagPath] = &componentInfo{
					Versions: []*semver.Version{componentVersion},
				}
			} else {
				existingVersion := semverVersions[tagPath]
				existingVersion.Versions = append(semverVersions[tagPath].Versions, componentVersion)
			}
			return nil
		})
	}
	versions := make(map[string]string)
	for k, v := range semverVersions {
		collection := semver.Collection(v.Versions)
		sort.Sort(collection)
		versions[k] = collection[len(collection)-1].String()
	}
	return versions
}
