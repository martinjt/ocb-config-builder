package main

import "github.com/martinjt/ocb-config-builder/pkg/configmapping"

func getCoreProcessorMapping(version string) map[string]configmapping.ComponentMapping {
	return map[string]configmapping.ComponentMapping{
		"batch": {
			GithubUrl: "go.opentelemetry.io/collector/processor/batchprocessor",
			Version:   version,
		},
	}
}

func getCoreExporterMapping(version string) map[string]configmapping.ComponentMapping {
	return map[string]configmapping.ComponentMapping{
		"logging": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/loggingexporter",
			Version:   version,
		},
		"otlp": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/otlpexporter",
			Version:   version,
		},
		"otlphttp": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/otlphttpexporter",
			Version:   version,
		},
	}
}

func getCoreReceiverMapping(version string) map[string]configmapping.ComponentMapping {
	return map[string]configmapping.ComponentMapping{
		"otlp": {
			GithubUrl: "go.opentelemetry.io/collector/receiver/otlpreceiver",
			Version:   version,
		},
	}
}
