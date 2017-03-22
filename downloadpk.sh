#!/bin/bash
if [ ! -f install.sh ]; then
echo 'install must be run within its container folder' 1>&2
exit 1
fi

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"


go get -d -v ./src/apidemo

export GOPATH="$OLDGOPATH"

echo 'finished'
