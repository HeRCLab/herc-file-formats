#!/bin/sh

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


# create a generic tarball
cd build
tar cvfz "../release/$RELEASENAME.tar.gz" .

# now create a deb package

mkdir ./DEBIAN
mkdir ./usr
mv ./bin ./include ./lib ./man ./usr/

touch ./DEBIAN/conffiles # none, just need an empty file her

echo "Package: herc-file-formats" > ./DEBIAN/control
echo "Priority: extra" >> ./DEBIAN/control
echo "Section: checkinstall" >> ./DEBIAN/control # I have no idea what "Section" means, just cargo-culted it from an example
echo "Installed-Size: $(du -csh --block-size=1kB . | tail -n1 | awk '{print($1)}')" >> ./DEBIAN/control
echo "Maintainer: None" >> ./DEBIAN/control
echo "Architecture: amd64" >> ./DEBIAN/control
echo "Version: $VERSION" >> ./DEBIAN/control
echo "Provides: herc-file-formats" >> ./DEBIAN/control
echo "Description: HeRC Lab tools and libraries for file formats" >> ./DEBIAN/control

dpkg-deb -b ./

mv '..deb' "../release/$RELEASENAME.deb"

rm -rf ./DEBIAN
rm -rf ./usr

cd ..
