#!/usr/bin/env bash

set -e

go fmt ./...

diff="$(git diff --stat -- '**/*.go')"

if [[ "$diff" != '' ]]; then
  echo
  echo The following files need to be reformatted:
  echo "$diff"
  echo
  echo Run \"go mod fmt ./...\"
  sleep 1
  exit 1
fi
