#!/bin/bash 
#
# This is run outside the container
#

ARCH="amd64"
# VERSION=${1}

set -Eeuo pipefail


die() {
  echo "$@"
  exit
}


cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1
BUILD_DIR=$(pwd)
cd ../../..
[ ! -f ".version" ] && die you need a .version file
source .version && echo "loaded"
COMPILE_DIR=$(pwd)
make build
cd $BUILD_DIR

PKG_DIR="${APPNAME}_v${VERSION}-1_${ARCH}"
mkdir ${PKG_DIR}
mkdir ${PKG_DIR}/DEBIAN
mkdir -p ${PKG_DIR}/usr/bin/
cp ../../../sqlchk ${PKG_DIR}/usr/bin/
cat <<EOT > ${PKG_DIR}/DEBIAN/control 
Package: ${APPNAME}
Version: ${VERSION}
Architecture: ${ARCH}
Maintainer: Joe Siwiak <info@sunset-crew.com>
Description: A program to query and load data into a db
EOT
export VERSION=${VERSION}
docker compose up 

# dpkg-deb --build --root-owner-group ${PKG_DIR}
