FROM local/simple-collector-builder:0.90.1-0.1 as build
ENV CONFIG_INPUT_LOCATION=/config.yaml
ENV ADDITIONAL_REPOS=github.com/martinjt/resourceattrtocontextconnector
COPY config-test.yaml /config.yaml
COPY config-test-additional.yaml /config-additional.yaml
RUN ocbconfigbuilder -d
RUN CGO_ENABLED=0 builder --config=/ocb-config.yaml --output-path=/app


FROM cgr.dev/chainguard/static:latest
COPY --from=build /app/otelcol-custom /
COPY --from=build /config.yaml /config.yaml
EXPOSE 4317/tcp 4318/tcp 13133/tcp

CMD ["/otelcol-custom", "--config=/config.yaml"]
