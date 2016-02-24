#!/bin/bash

# This script downloads and installs the NGINX load balancer.
# Execute using:
#
#     source ./nginx_install {install directory} {version} {openssl install directory} {number of cores (optional)}

INSTALL_ROOT=$1
NGINX_VERSION=$2
OPENSSL_PATH=$3
NUM_CORES=$4

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$NGINX_VERSION" ]]; then
	echo "No NGINX version specified."
	exit 1
fi

if [[ -z "$NUM_CORES" ]]; then
	NUM_CORES="2"
fi

# Variables
SOURCE_DIR=$INSTALL_ROOT/nginx-$NGINX_VERSION
INSTALL_DIR=$INSTALL_ROOT/_nginx_

echo
echo "Downloading and installing the NGINX load balancer..."
echo

# Save current directory
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download source
wget http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz
tar -xvf nginx-$NGINX_VERSION.tar.gz
rm -f nginx-$NGINX_VERSION.tar.gz

# Build and install
cd $SOURCE_DIR
./configure --prefix=$INSTALL_DIR --without-http_rewrite_module --with-openssl=$OPENSSL_PATH
make
make install

# Symlink executables to install root
ln -s $INSTALL_DIR/bin/* $INSTALL_ROOT/

# Clean up
cd $CURRENT_DIR
rm -rf $SOURCE_DIR

echo
echo "NGINX installation complete."
