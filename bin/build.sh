#!/usr/bin/env bash

GREEN='\033[1;92m'
NC='\033[0m'
print() {
	echo -e "${GREEN}> -===[ ${1} ]===- ${NC}"
}

rm -rf dist
mkdir -p dist || exit 1
print 'frontend build'
cd frontend
yarn
yarn build
EXITCODE=$?
cd ..
if [ $EXITCODE != 0 ]; then
   exit $EXITCODE
fi

if [ "$GOOS" == "windows" ]; then
	EXE_SUFFIX=".exe"
else
	EXE_SUFFIX=""
fi

set -e
print 'go build hanectl'
env CGO_ENABLED=0 $GOROOT/bin/go build -trimpath -a -ldflags '-s -w -extldflags "-static"' \
  -o dist/hanectl${EXE_SUFFIX} hanectl/main.go
#env CGO_ENABLED=1 $GOROOT/bin/go build -trimpath -a -ldflags '-s -w -extldflags' -o dist/hanectl${EXE_SUFFIX} hanectl/main.go

print 'go build password'
env CGO_ENABLED=0 $GOROOT/bin/go build -trimpath -a -ldflags '-s -w -extldflags "-static"' \
  -o dist/hanectl_pwdgen${EXE_SUFFIX} hanectl/cmd/password/main.go


print 'build dist'
cp -r config dist/
cp -rf frontend/build dist/web

#if test -f "dist/config/scripts/_lib.min.js"; then
#  rm -f dist/config/scripts/_lib.js
#  mv dist/config/scripts/_lib.min.js dist/config/scripts/_lib.js
#fi
