#!/bin/sh

echo "This script bricked my laptop the last time I ran it, and I have not figured out why yet."
echo "You should not run this script unless you really know what you're doing."
echo "-- Charles"
exit 1

set -e
set -u
set -x

# Generates a collection of release tarballs

ARCH="$(uname -p)"

cd "$(dirname "$0")"


make clean
make build

mkdir -p ./release

VERSION="$(cat ./VERSION)"

RELEASENAME="herc-file-formats-$VERSION-$ARCH"

cd build
tar cvfz "../release/$RELEASENAME.tar.gz" .
cd ..

PROJ="$(pwd)"
TEMP="$(mktemp -d)"

cp -R ./ "$TEMP"
cd "$TEMP"

printf "HeRC File Formats\n" | sudo checkinstall -D --install=no --gzman --strip --nodoc --pkgrelease "$VERSION" --pkgname herc-file-formats
sudo chown "$(whoami)" *.deb
mv *.deb "$PROJ/release/$RELEASENAME.deb"

ls -lah

cd "$PROJ"
rm -rf "$TEMP"
