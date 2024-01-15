#!/bin/bash 
#
#
#

PKG="sqlchk"
ARCH="amd64"
VERSION=${1}

set -Eeuo pipefail


die() {
  echo "$@"
  exit
}


cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1
BUILD_DIR=$(pwd)
cd ../..
[ ! -f ".version" ] && die you need a .version file
source .version && echo "loaded"
cd ${BUILD_DIR}

PKG_DIR="${APPNAME}_v${VERSION}-1_$ARCH"
echo $PKG_DIR
if [ -d "$PKG_DIR" ]; then
  rm -rf $PKG_DIR
  rm -rf $PKG_DIR.deb
  echo "deleteable"
fi


