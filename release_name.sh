#!/bin/bash

# example tag name 1.0-collector-0.104.0-release
regex='^([0-9]+\.[0-9]+)-collector-([0-9]+\.[0-9]+\.[0-9]+)-release$'

if [[ $1 =~ $regex ]]; then
  echo "BUILDER_VERSION=${BASH_REMATCH[1]}"
  echo "COLLECTOR_VERSION=${BASH_REMATCH[2]}"
else
  echo "Invalid tag name"
fi
