#!/bin/bash

# This script downloads and installs the Mono Framework from source.
# Execute using:
#
#     source ./mono_install {install directory} {version} {number of available cpu cores (optional)}

INSTALL_ROOT=$1
MONO_VERSION=$2
MINOR_RELEASE=$3
NUM_CORES=$4

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$MONO_VERSION" ]]; then
	echo "No Mono version specified."
	exit 1
fi

if [[ -z "$MINOR_RELEASE" ]]; then
	echo "No minor release number specified."
	exit 1
fi

if [[ -z "$NUM_CORES" ]]; then
	NUM_CORES="2"
fi

# Variables
SOURCE_DIR=$INSTALL_ROOT/mono-$MONO_VERSION
INSTALL_DIR=$INSTALL_ROOT/_mono_

echo
echo "Downloading and installing Mono Framework..."
echo

# Save current directory
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download source
wget http://download.mono-project.com/sources/mono/mono-$MONO_VERSION.$MINOR_RELEASE.tar.bz2
tar -xvf mono-$MONO_VERSION.$MINOR_RELEASE.tar.bz2
rm -f mono-$MONO_VERSION.$MINOR_RELEASE.tar.bz2

# Build and install
cd $SOURCE_DIR
./configure --prefix=$INSTALL_DIR
make
make install

# Symlink to install directory
ln -s $INSTALL_DIR/bin/mono $INSTALL_ROOT/mono
ln -s $INSTALL_DIR/bin/xbuild $INSTALL_ROOT/xbuild

# Clean up
cd $CURRENT_DIR
rm -rf $SOURCE_DIR


echo
echo "Mono Framework installation complete."
