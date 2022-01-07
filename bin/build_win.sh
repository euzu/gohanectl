#!/usr/bin/env bash
set -e
DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
env GOOS=windows GOARCH=amd64  $DIR/build.sh
exit $?
