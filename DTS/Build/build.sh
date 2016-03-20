#!/bin/bash

REPO_ROOT=$1

if [[ -z "$REPO_ROOT" ]]; then
	echo "Repository root not set."
	exit 1
fi

cd $REPO_ROOT

# Clean bin directory
rm -rf $REPO_ROOT/bin
mkdir $REPO_ROOT/bin

# Build DTS Services
DTS/Build/NuGet/NuGet.exe restore DTS/Services/DTSServices.sln
xbuild /p:Configuration=Release DTS/Services/DTSServices.sln

# Build API 
export GOPATH="$REPO_ROOT/DTS/API"
mkdir $REPO_ROOT/bin/DtsApi
cd $REPO_ROOT/DTS/API/src/dtsapi
go get
go install
mv $REPO_ROOT/DTS/API/bin/dtsapi $REPO_ROOT/bin/DtsApi/DtsApi
cp $REPO_ROOT/DTS/API/src/dtsapi/conf.json $REPO_ROOT/bin/DtsApi
rm -rf $REPO_ROOT/DTS/API/bin

# Build QuoteCache
export GOPATH="$REPO_ROOT/DTS/QuoteCache"
mkdir $REPO_ROOT/bin/QuoteCache
cd $REPO_ROOT/DTS/QuoteCache/src/quotecache
go get
go install
mv $REPO_ROOT/DTS/QuoteCache/bin/quotecache $REPO_ROOT/bin/QuoteCache/QuoteCache
cp $REPO_ROOT/DTS/QuoteCache/src/quotecache/conf.json $REPO_ROOT/bin/QuoteCache
rm -rf $REPO_ROOT/DTS/QuoteCache/bin
