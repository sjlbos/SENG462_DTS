#!/bin/bash

INSTALL_ROOT=$1

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

cd $INSTALL_ROOT
# Enable management console
./rabbitmq-plugins enable rabbitmq_management

# Start service
./rabbitmq-server -detached

# Add and Configure User
./rabbitmqctl add_user dts_user Group1

./rabbitmqctl set_permissions -p / dts_user "^dts_user-.*" ".*" ".*"

