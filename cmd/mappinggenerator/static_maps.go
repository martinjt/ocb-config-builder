package main

import "github.com/martinjt/ocb-config-builder/pkg/configmapping"

func getCoreProcessorMapping(version string) map[string]configmapping.ComponentMapping {
	return map[string]configmapping.ComponentMapping{
		"batch": {
			GithubUrl: "go.opentelemetry.io/collector/processor/batchprocessor",
			Version:   version,
		},
		"memory_limiter": {
			GithubUrl: "go.opentelemetry.io/collector/processor/memorylimiterprocessor",
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
		"debug": {
			GithubUrl: "go.opentelemetry.io/collector/exporter/debugexporter",
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

func getCoreExtensionMapping(version string) map[string]configmapping.ComponentMapping {
	return map[string]configmapping.ComponentMapping{
		"memory_ballast": {
			GithubUrl: "go.opentelemetry.io/collector/extension/memoryballastextension",
			Version:   version,
		},
		"memory_limiter": {
			GithubUrl: "go.opentelemetry.io/collector/extension/memorylimiterextension",
			Version:   version,
		},
		"zpages": {
			GithubUrl: "go.opentelemetry.io/collector/extension/zpagesextension",
			Version:   version,
		},
	}
}

func getCoreConnectorMapping(version string) map[string]configmapping.ComponentMapping {
	return map[string]configmapping.ComponentMapping{
		"forward": {
			GithubUrl: "go.opentelemetry.io/collector/connector/forwardconnector",
			Version:   version,
		},
	}
}
