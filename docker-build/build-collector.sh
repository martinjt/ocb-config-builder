#!/bin/sh

/builder/ocbconfigbuilder --input $1 --output /tmp/ocb_config.yaml
CGO_ENABLED=0 builder --config=/tmp/ocb_config.yaml --output-path=/app