#!/bin/bash

# This script downloads and installs the Erlang runtime from source.
# Execute using:
#
#    	./erlang_install {install directory} {version} {number of available cpu cores (optional)}

INSTALL_ROOT=$1
ERLANG_VERSION=$2
NUM_CORES=$3

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$ERLANG_VERSION" ]]; then
	echo "No Erlang version specified."
	exit 1
fi

if [[ -z "$NUM_CORES" ]]; then
	NUM_CORES="2"
fi

echo
echo "Downloading and installing erlang..."
echo

CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download source
wget http://www.erlang.org/download/otp_src_$ERLANG_VERSION.tar.gz
tar -zxvf otp_src_$ERLANG_VERSION.tar.gz
rm -f otp_src_$ERLANG_VERSION.tar.gz

# Build and install
SOURCE_DIR=$INSTALL_ROOT/otp_src_$ERLANG_VERSION
INSTALL_DIR=$INSTALL_ROOT/erlang/$ERLANG_VERSION

cd $SOURCE_DIR
export ERL_TOP=$SOURCE_DIR
./configure prefix=$INSTALL_DIR
make -j $NUM_CORES
make install

# Symlink to install directory
ln -s $INSTALL_DIR/bin/erl $INSTALL_ROOT/erl

# Clean up
cd ..
rm -rf $SOURCE_DIR
cd $CURRENT_DIR

echo
echo "Erlang installation complete."
