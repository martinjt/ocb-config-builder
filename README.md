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
COPY config-test.yaml /config.yaml
RUN ocbconfigbuilder
RUN CGO_ENABLED=0 builder --config=/ocb-config.yaml --output-path=/app

FROM cgr.dev/chainguard/static:latest
COPY --from=build /app/otelcol-custom /
COPY config.yaml /
EXPOSE 4317/tcp 4318/tcp 13133/tcp

CMD ["/otelcol-custom", "--config=/config.yaml"]
```

This will result in a custom image that has the config file embedded inside, however, there is no reason that you can't mount another config file with the same processors in the same place, or not provide the config file in the runtime image at all (to force it to be mounted later).  

## Building with custom components

If you would like to include components that are not in the in the otel collector Core and Contrib repos, this is possible using the builder, however there are requirement on the external repo.

### metadata.yaml

The repo where the component exists must include a `metadata.yaml` file in the repo. This should look the same as those in contrib, specifically, they should look like this:

```yaml
type: resourceattr_to_context

status:
  class: connector
  stability:
    alpha: [traces_to_metrics]
  distributions: []
  codeowners:
    active: [martinjt]
```

The important parts are the `type` attribute, which must match the attribute name in the collector config that will be used, and the `status::class` attribute which must match to the type of component this is.

Additionally, if there are other `metadata.yaml` files in the repo, they must follow this format.

### Git tags

The repository must also have a tag for the component that must be in the format `{path-to-component}/v{version}`. This is the same as the OCB itself since we're using the same methods.

`{path-to-component}` is the directory path from the repository root.

The builder will pull all the tags, and choose the latest version as the one to use for the config.

### ADDITIONAL_REPOS environment variable

The `ocbconfigbuilder` application looks for a comma separated list of repository url's to search for the metadata. This should work with any git server that allow anonymous authentication and allows http communication.

If you're unsure, you can do the following on a machine that is not authenticated to that repo.

```shell
git clone https://{my-repo}.git
```

If this works, then if you add the same string as `{my-repo}` to your `ADDITIONAL_REPOS` environment variable it should work.

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