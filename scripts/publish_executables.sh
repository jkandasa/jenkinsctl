#!/bin/bash

# this script used to generate binary files
# should be executed from the root locations of the repository


source ./scripts/version.sh

BUILD_DIR=builds
BINARY_DIR=binary
# clean builds directory
rm ${BUILD_DIR}/* -rf

# create directories
mkdir -p ${BUILD_DIR}/${BINARY_DIR}

# download dependencies
go mod tidy

function package {
  local PACKAGE_STAGING_DIR=$1
  local BINARY_FILE=$2
  local FILE_EXTENSION=$3

  mkdir -p ${PACKAGE_STAGING_DIR}

  # echo "Package dir: ${PACKAGE_STAGING_DIR}"
  cp ${BUILD_DIR}/${BINARY_DIR}/${BINARY_FILE} ${PACKAGE_STAGING_DIR}/jenkinsctl${FILE_EXTENSION}

  # copy license
  cp LICENSE ${PACKAGE_STAGING_DIR}/LICENSE.txt

  ARCHIVE_NAME="${PACKAGE_STAGING_DIR}.tar.gz"
  # echo "Packaging into: ${ARCHIVE_NAME}"
  tar -czf ${BUILD_DIR}/${ARCHIVE_NAME} ${PACKAGE_STAGING_DIR}
  rm ${PACKAGE_STAGING_DIR} -rf
}


# platforms to build
PLATFORMS=("linux/arm" "linux/arm64" "linux/386" "linux/amd64" "linux/ppc64" "linux/ppc64le" "linux/s390x" "darwin/amd64" "windows/386" "windows/amd64")

# compile
for platform in "${PLATFORMS[@]}"
do
  platform_raw=(${platform//\// })
  GOOS=${platform_raw[0]}
  GOARCH=${platform_raw[1]}
  package_name="jenkinsctl-${GOOS}-${GOARCH}"

  env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${BUILD_DIR}/${BINARY_DIR}/${package_name} -ldflags "$LD_FLAGS" cmd/main.go
  if [ $? -ne 0 ]; then
    echo 'an error has occurred. aborting the build process'
    exit 1
  fi

  FILE_EXTENSION=""
  if [ $GOOS = "windows" ]; then
    FILE_EXTENSION='.exe'
  fi

  package jenkinsctl-${VERSION}-${GOOS}-${GOARCH} ${package_name} ${FILE_EXTENSION}

done
