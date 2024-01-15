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

PKG_DIR="${PKG}_v${VERSION}-1_$ARCH"

# docker compose up

dpkg-deb --build --root-owner-group ${PKG_DIR}
