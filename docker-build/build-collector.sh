#!/bin/sh

/builder/ocbconfigbuilder --input $1 --output /tmp/ocb_config.yaml
builder --config=/tmp/ocb_config.yaml --output-path=/app
#cat /tmp/ocb_config.yaml