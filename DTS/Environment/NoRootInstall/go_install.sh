#!/bin/bash

# This script downloads and installs the Go programming language runtime.
# Execute using:
#
#    	./go_install {install directory} {go version}

INSTALL_ROOT=$1
GO_VERSION=$2

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$GO_VERSION" ]]; then
	GO_VERSION="1.5.3"
fi

echo
echo "Downloading and installing Go..."
echo

INSTALL_DIR= $INSTALL_ROOT/go_src
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download and extract
wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz
tar -xzf go$GO_VERSION.linux-amd64.tar.gz
rm -f go$GO_VERSION.linux-amd64.tar.gz

# Symlink to install directory
mv go go_src

ln -s ./go_src/bin/go ./go
ln -s ./go_src/bin/godoc ./godoc
ln -s ./go_src/bin/gofmt ./goftm

cd $CURRENT_DIR

echo
echo "Go installation complete..."
