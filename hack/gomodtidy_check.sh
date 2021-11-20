#!/usr/bin/env bash

go mod tidy

diff="$(git diff -- 'go.mod' 'go.sum')"

if [[ "$diff" != '' ]]; then
  echo
  echo "'go mod tidy' needs to be run again"
  echo
  echo "$diff"
  sleep 1
  exit 1
fi
