#!/bin/bash

# This script downloads and installs a Postgresql database server from source.
# Execute using:
#
#     source ./postgres_install {install directory} {version} {number of available cpu cores (optional)}

INSTALL_ROOT=$1
PSQL_VERSION=$2
NUM_CORES=$3

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$PSQL_VERSION" ]]; then
	echo "No Postgresql version specified."
	exit 1
fi

if [[ -z "$NUM_CORES" ]]; then
	NUM_CORES="2"
fi

# Variables
SOURCE_DIR=$INSTALL_ROOT/postgresql-$PSQL_VERSION
INSTALL_DIR=$INSTALL_ROOT/postgresql

echo
echo "Downloading and installing Postgresql..."
echo

# Save current directory
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download source
wget https://ftp.postgresql.org/pub/source/v$PSQL_VERSION/postgresql-$PSQL_VERSION.tar.gz
tar -zxvf postgresql-$PSQL_VERSION.tar.gz
rm -f postgresql-$PSQL_VERSION.tar.gz
mv $INSTALL_ROOT/postgres $INSTALL_DIR

# Build and install
cd $SOURCE_DIR
./configure --prefix=$INSTALL_DIR
make -j $NUM_CORES
make install

# Symlink executables to install root
ln -s $INSTALL_DIR/bin/* $INSTALL_ROOT/

# Clean up
cd ..
rm -rf $SOURCE_DIR
cd $CURRENT_DIR

echo
echo "Postgresql installation complete."
