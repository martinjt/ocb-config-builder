FROM local/simple-collector-builder:0.90.1 as build
COPY config-test.yaml /config/config.yaml
RUN /builder/build-collector.sh /config/config.yaml

FROM cgr.dev/chainguard/static:latest
COPY --from=build /app/otelcol-custom /
COPY config-test.yaml /
EXPOSE 4317/tcp 4318/tcp 13133/tcp

CMD ["/otelcol-custom", "--config=/config-test.yaml"]
