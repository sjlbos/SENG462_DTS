#!/bin/bash

# This script downloads and installs the MongoDB database.
# Execute using:
#
#     source ./mongo_install {install directory} {version}

INSTALL_ROOT=$1
MONGO_VERSION=$2

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$MONGO_VERSION" ]]; then
	echo "No Mono version specified."
	exit 1
fi

# Variables
TARGET_PLATFORM="mongodb-linux-x86_64-rhel70"
INSTALL_DIR=$INSTALL_ROOT/$TARGET_PLATFORM-$MONGO_VERSION

echo
echo "Downloading and installing MongoDB..."
echo

# Save current directory
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download source
wget https://fastdl.mongodb.org/linux/$TARGET_PLATFORM-$MONGO_VERSION.tgz
tar -zxvf $TARGET_PLATFORM-$MONGO_VERSION.tgz
rm -f $TARGET_PLATFORM-$MONGO_VERSION.tgz

# Symlink to install directory
ln -s $INSTALL_DIR/bin/* $INSTALL_ROOT/

# Clean up
cd $CURRENT_DIR

echo
echo "MongoDB installation complete."
