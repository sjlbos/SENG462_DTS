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
	echo "No Go version specified."
	exit 1
fi

echo
echo "Downloading and installing Go..."
echo

INSTALL_DIR = $INSTALL_ROOT/go
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download and extract
wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz
tar -xzf go$GO_VERSION.linux-amd64.tar.gz
rm -f go$GO_VERSION.linux-amd64.tar.gz

# Symlink to install directory
ln -s $INSTALL_DIR/bin/go $INSTALL_ROOT/go
ln -s $INSTALL_DIR/bin/godoc $INSTALL_ROOT/godoc
ln -s $INSTALL_DIR/bin/gofmt $INSTALL_ROOT/goftm

cd $CURRENT_DIR

echo
echo "Go installation complete..."