#!/bin/bash

URL=$(curl -s https://api.github.com/repos/ip2location/ip2location-io-cli/releases/latest | grep 'download_url.*darwin_amd64' | awk '{ print $2 }' | sed 's/"//g')
FILE=$(echo $URL | xargs basename)

curl -LO $URL
tar -xf $FILE
rm $FILE
mv $(echo $FILE | sed 's/.tar.gz$//g') /usr/local/bin/ip2locationio

echo
echo 'You can now run `ip2locationio`'.

if [ -f "$0" ]; then
	rm $0
fi