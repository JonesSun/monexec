#! /bin/bash

if [ ! -n ${1} ];then
    echo "please input GOOS"
    exit
fi
if [ ! -n ${2} ];then
    echo "please input GOARCH"
    exit
fi

export GOOS=${1}
export GOARCH=${2}
export GOPATH=`pwd`

if [ "windows" == ${GOOS} ];then
    export fileName=monexec.exe
   else
    export fileName=monexec
fi
go build -o bin/${fileName} -ldflags "-w -s" github.com/reddec/monexec
upx -9 bin/${fileName}