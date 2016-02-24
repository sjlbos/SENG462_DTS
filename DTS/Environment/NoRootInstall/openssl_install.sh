#!/bin/bash

# This script downloads and installs OpenSSL.
# Execute using:
#
#     source ./openssl_install {install directory} {version} {number of cores (optional)}

INSTALL_ROOT=$1
OPENSSL_VERSION=$2
NUM_CORES=$3

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$OPENSSL_VERSION" ]]; then
	echo "No OpenSSL version specified."
	exit 1
fi

if [[ -z "$NUM_CORES" ]]; then
	NUM_CORES="2"
fi

# Variables
SOURCE_DIR=$INSTALL_ROOT/openssl-$OPENSSL_VERSION
INSTALL_DIR=$INSTALL_ROOT/_openssl_

echo
echo "Downloading and installing OpenSSL..."
echo

# Save current directory
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

mkdir -p $INSTALL_DIR

# Download source
wget http://www.openssl.org/source/openssl-$OPENSSL_VERSION.tar.gz
tar -xvf openssl-$OPENSSL_VERSION.tar.gz
rm -f openssl-$OPENSSL_VERSION.tar.gz

# Build and install
cd $SOURCE_DIR
./config --prefix=$INSTALL_DIR no-shared no-zlib -fPIC no-gost
make depend
make -j $NUM_CORES
make install

# Symlink executables to install root
ln -s $INSTALL_DIR/bin/* $INSTALL_ROOT/

# Clean up
cd $CURRENT_DIR
rm -rf $SOURCE_DIR

echo
echo "OpenSSL installation complete."
