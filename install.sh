#!/bin/bash
if [ ! -f install.sh ]; then
echo 'install must be run within its container folder' 1>&2
exit 1
fi

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

#gofmt -w src
echo $GOPATH
go install apidemo 

#copy the config file  to bin
cp $CURDIR/src/config/logconfig.xml ./bin

export GOPATH="$OLDGOPATH"

echo 'finished'
