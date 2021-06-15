#!/usr/bin/env bash
set -e
DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
env GOOS=linux GOARCH=amd64 "$DIR"/build.sh
exit $?
