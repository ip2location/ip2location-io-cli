#!/bin/bash

# Create the downloads

set -e

DIR=`dirname $0`
ROOT=$DIR/..

NAME=$1
VERSION=$2

# Create Go binaries
rm -f $ROOT/build/${NAME}_${VERSION}*
rm -f $ROOT/build/ip2location-io-cli_*.deb
$ROOT/scripts/build-all.sh "${NAME}" "$VERSION"

# archive
cd $ROOT/build
for t in ${NAME}_${VERSION}_* ; do
    if [[ $t == ${NAME}_*_windows_* ]]; then
        zip -q ${t/.exe/.zip} $t
    else
        tar -czf ${t}.tar.gz $t
    fi
done
cd ..

# Create Deb package
rm -rf $ROOT/dist
mkdir -p $ROOT/dist/DEBIAN/
mkdir -p $ROOT/dist/usr/local/bin
cp $ROOT/debian/control $ROOT/dist/DEBIAN/
cp $ROOT/build/${NAME}_${VERSION}_linux_amd64 $ROOT/dist/usr/local/bin/${NAME}
dpkg-deb -Zgzip --build ${ROOT}/dist build/ip2location-io-${VERSION}.deb

# Create architecture-specific Deb packages
for arch in amd64 arm64 ; do
    rm -rf $ROOT/dist
    mkdir -p $ROOT/dist/DEBIAN/
    mkdir -p $ROOT/dist/usr/local/bin

    # Stamp the package name, version and architecture into the control file
    sed -e "s/^Package: .*/Package: ip2location-io-cli/"      \
        -e "s/^Architecture: .*/Architecture: ${arch}/"       \
        -e "s/^Version: .*/Version: ${VERSION}/"              \
        $ROOT/debian/control > $ROOT/dist/DEBIAN/control

    cp $ROOT/build/${NAME}_${VERSION}_linux_${arch} $ROOT/dist/usr/local/bin/${NAME}
    dpkg-deb -Zgzip --build ${ROOT}/dist ${ROOT}/build/ip2location-io-cli_${arch}.deb
done

rm -rf $ROOT/dist
