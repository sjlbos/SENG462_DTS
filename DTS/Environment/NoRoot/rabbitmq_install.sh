#!/bin/bash

# This script downloads and installs a RabbitMQ broker server.
# Execute using:
#
#     source ./rabbitmq_install {install directory} {version}

INSTALL_ROOT=$1
RABBIT_VERSION=$2

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "$RABBIT_VERSION" ]]; then
	echo "No RabbitMQ version specified."
	exit 1
fi

# Variables
INSTALL_DIR=$INSTALL_ROOT/rabbitmq_server-$RABBIT_VERSION

echo
echo "Downloading and installing RabbitMQ..."
echo

# Save current directory
CURRENT_DIR=`pwd`

mkdir -p $INSTALL_ROOT
cd $INSTALL_ROOT

# Download binaries
wget https://www.rabbitmq.com/releases/rabbitmq-server/v$RABBIT_VERSION/rabbitmq-server-generic-unix-$RABBIT_VERSION.tar.xz
tar -xJf rabbitmq-server-generic-unix-$RABBIT_VERSION.tar.xz
rm -f rabbitmq-server-generic-unix-$RABBIT_VERSION.tar.xz

# Add RabbitMQ to the PATH
export PATH=$PATH:$INSTALL_DIR/sbin

# Enable management console and start service
rabbitmq-plugins enable rabbitmq_management
rabbitmq_server-3.6.0/sbin/rabbitmq-server -detached

# Clean up
cd $CURRENT_DIR

echo
echo "RabbitMQ installation complete."
