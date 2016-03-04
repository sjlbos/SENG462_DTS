#!/bin/bash

REPO_ROOT=$1

if [[ -z "$REPO_ROOT" ]]; then
	echo "Repository root not set."
	exit 1
fi

cd $REPO_ROOT

# Build DTS Services
xbuild /p:Configuration=Release DTS/Services/DTSServices.sln

# Set temporary GOPATH
GOPATH="$REPO_ROOT/DTS/API:$REPO_ROOT/DTS/QuoteCache"

# Build API 
export GOPATH="$REPO_ROOT/DTS/API"
cd $REPO_ROOT/DTS/API/src/dtsapi
go get
go install
mv $REPO_ROOT/DTS/API/bin/dtsapi $REPO_ROOT/bin/DtsApi
rm -rf $REPO_ROOT/DTS/API/bin

# Build QuoteCache
export GOPATH="$REPO_ROOT/DTS/QuoteCache"
cd $REPO_ROOT/DTS/QuoteCache/src/quotecache
go get
go install
mv $REPO_ROOT/DTS/QuoteCache/bin/quotecache $REPO_ROOT/bin/QuoteCache
rm -rf $REPO_ROOT/DTS/QuoteCache/bin









