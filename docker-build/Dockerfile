FROM cgr.dev/chainguard/go:latest as build
ARG OTEL_VERSION=0.90.1
ARG CONFIG_OUTPUT_LOCATION=/ocb-config.yaml
ARG CONFIG_INPUT_LOCATION=/config.yaml
ENV CONFIG_INPUT_LOCATION=${CONFIG_INPUT_LOCATION}
ENV CONFIG_OUTPUT_LOCATION=${CONFIG_OUTPUT_LOCATION}
ENV GOPATH=/usr
WORKDIR /usr/bin
RUN go install go.opentelemetry.io/collector/cmd/builder@v${OTEL_VERSION}
COPY ocbconfigbuilder .
RUN chmod +x ocbconfigbuilder
