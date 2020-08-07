#!/bin/bash

if [ -z $TARGET_HOST ]
then
	echo "Building uplink-c ..."
	go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-archive -tags stdsha256 -o libuplinkc.a .
elif [ $TARGET_HOST -a "x86_64-w64-mingw32" ]; then 
	HOST_OS=`go env | grep "GOHOSTOS" | grep -o '"[^"]\+"' | sed 's/"//g'`
	if [ $HOST_OS == "linux" ]
	then 
		echo "Cross-compiling uplink-c for Windows ..."
		GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CXX="x86_64-w64-mingw32-g++" CC="x86_64-w64-mingw32-gcc" go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-archive -tags stdsha256 -o libuplinkc.a .
	fi
fi