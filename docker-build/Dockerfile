FROM cgr.dev/chainguard/go:latest as build
ARG  OTEL_VERSION=0.90.1
ENV GOPATH=/usr
WORKDIR /builder
RUN go install go.opentelemetry.io/collector/cmd/builder@v${OTEL_VERSION}
COPY ocbconfigbuilder .
COPY build-collector.sh .
RUN chmod +x build-collector.sh

CMD ["/builder/build-collector.sh", "/config/config.yaml" ]