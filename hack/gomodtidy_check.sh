#!/usr/bin/env

go mod tidy

if [[ $(git diff --stat) != '' ]]; then
  echo
  echo "'go mod tidy' needs to be run again"
  echo
  sleep 1
  exit 1
fi
