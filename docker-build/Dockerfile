FROM golang:1.21 as build
ARG  OTEL_VERSION=0.90.1
WORKDIR /builder
RUN go install go.opentelemetry.io/collector/cmd/builder@v${OTEL_VERSION}
COPY ocbconfigbuilder .
COPY build-collector.sh .
RUN chmod +x build-collector.sh

CMD ["/builder/build-collector.sh", "/config/config.yaml" ]