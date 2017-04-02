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

#echo $0 $1
if [ $1="debug" ]
then
    go install client 
else
    echo "will build to release"   
    #go install -ldflags “-s” client  
    go install client  
fi





export GOPATH="$OLDGOPATH"

echo 'finished'
