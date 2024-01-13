#!/bin/bash
#
#
#

set -Eeuo pipefail


die() {
  echo "$@"
  exit
}


cd "$(dirname "${BASH_SOURCE[0]}")"/../.. >/dev/null 2>&1

[ ! -f ".version" ] && die you need a .version file 
source .version && echo "loaded" || die "it didn't load $(ls)"

mkdir deploy/rpm/build
mkdir deploy/rpm/build/BUILD
mkdir deploy/rpm/build/RPMS
mkdir deploy/rpm/build/SRPMS
mkdir deploy/rpm/build/SOURCES
mkdir deploy/rpm/build/SPEC

cp deploy/rpm/${APPNAME}.spec deploy/rpm/build/SPEC/${APPNAME}.spec
rsync -ar --exclude 'deploy' --exclude '.git' --exclude 'notes' . ${APPNAME}-v${VERSION}/
tar czf deploy/rpm/build/SOURCES/${APPNAME}-v${VERSION}.tar.gz ${APPNAME}-v${VERSION}/
rm -rf ${APPNAME}-v${VERSION}/
cd deploy/rpm/
# docker compose up --build
echo "building ${VERSION}"
echo "VERSION=${VERSION}" > .env
docker compose up --build
[ -f ".env" ] && rm -vf .env
