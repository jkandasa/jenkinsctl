#!/bin/bash

# this script used to generate binary files
# should be executed from the root locations of the repository


source ./scripts/version.sh

BUILD_DIR=builds
# clean builds directory
rm ${BUILD_DIR}/* -rf


# create directories
mkdir -p ${BUILD_DIR}

# download dependencies
go mod tidy


# platforms to build
PLATFORMS=("linux/arm" "linux/arm64" "linux/386" "linux/amd64" "linux/ppc64" "linux/ppc64le" "linux/s390x" "darwin/386" "darwin/amd64" "windows/386" "windows/amd64")

# compile
for platform in "${PLATFORMS[@]}"
do
  platform_raw=(${platform//\// })
  GOOS=${platform_raw[0]}
  GOARCH=${platform_raw[1]}

  FILE_EXTENSION=""
  if [ $GOOS = "windows" ]; then
    FILE_EXTENSION='.exe'
  fi

  package_name="jenkinsctl-${VERSION}-${GOOS}-${GOARCH}${FILE_EXTENSION}"

  env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${BUILD_DIR}/${package_name} -ldflags "$LD_FLAGS" cmd/main.go
  if [ $? -ne 0 ]; then
    echo 'an error has occurred. aborting the build process'
    exit 1
  fi

done

