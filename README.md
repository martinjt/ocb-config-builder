# OpenTelemetry Collector Builder for Collector Config Files

This image allows you to mount a Standard OpenTelemetry config file and it will build a custom collector with only those specific components enabled.

This is a security best practice for the collector, but has been hard to work out due to the complexity involved in the standard OpenTelemetry Collector builder.

On top of the collector running a reduced component set (which reduces the potential vulnerability surface area), it also uses the Chainguard image for the build image, and the examples show how to use the Chainguard base images for running.

## Process

There are 2 types of components

* Core
* Contrib

For the core components, these need to be manually added in [mappinggenerator/main.go].

For contrib components, the mapping generator downloads the latest from the contrib repo, scans for the `metadata.yaml` files and adds any it finds in the exporter, receiver, extensions, processor, or connector directories.

This mapping file is embedded in the main `ocbconfigbuilder` binary.

The builder itself takes the names of all the components, and builds an OCB compatible file from it.

## Container image

There is an image you can use that means you do not need to know about any of the mappings.

A build will look like this

```dockerfile
FROM ghcr.io/martinjt/ocb-config-builder:latest as build
COPY config.yaml /config/config.yaml
RUN /builder/build-collector.sh /config/config.yaml

FROM cgr.dev/chainguard/static:latest
COPY --from=build /app/otelcol-custom /
COPY config.yaml /
EXPOSE 4317/tcp 4318/tcp 13133/tcp

CMD ["/otelcol-custom", "--config=/config.yaml"]
```

This will result in a custom image that has the config file embedded inside, however, there is no reason that you can't mount another config file with the same processors in the same place, or not provide the config file in the runtime image at all (to force it to be mounted later).  

## Building

To regenerate the mappings

```shell
go run cmd/mappinggenerator/*
```

To build the configbuilder

```shell
go build -o docker-build/ocbconfigbuilder cmd/configgenerator/*.go
```

To build the container image locally

```shell
docker build -f docker-build/Dockerfile -t local/simple-collector-builder:0.90.1 docker-build/
```

To build a custom collector image from the test config

```shell
docker build -f test/Dockerfile -t local\custom-collector:0.90.1 test/
```

## Current constaints

* It will not work with custom components
* Only some Core components are supported (simple to add more, just needs to be manual)
* Only works with Contrib components that have an upto data `metadata.yaml` (if you get a bug with a missing component, check that metadata file in the contrib repo)