FROM cgr.dev/chainguard/go:latest as build
ARG COLLECTOR_VERSION=0.92.0
ENV CONFIG_INPUT_LOCATION=/config.yaml
ENV CONFIG_OUTPUT_LOCATION=/ocb-config.yaml
ENV COLLECTOR_VERSION=${COLLECTOR_VERSION}
ENV GOPATH=/usr
WORKDIR /usr/bin
RUN go install go.opentelemetry.io/collector/cmd/builder@v${COLLECTOR_VERSION}
COPY ocbconfigbuilder .
RUN chmod +x ocbconfigbuilder
