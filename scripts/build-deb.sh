#!/bin/bash

VERSION="1.1.0"

rm -rf ../dist
mkdir -p ../dist/DEBIAN/
mkdir -p ../dist/usr/local/bin/
cp ../debian/control ../dist/DEBIAN/
cd ../ip2locationio
rm -f go.*
go mod init ip2location-io && go mod tidy && go build -o ../dist/usr/local/bin/ip2locationio
cd ..
dpkg-deb -Zgzip --build dist ip2location-io-$VERSION.deb