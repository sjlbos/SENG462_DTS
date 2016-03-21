#!/bin/bash

INSTALL_ROOT=$1
NUM_CORES=$2

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "No install root specified."
	exit 1
fi

if [[ -z "NUM_CORES" ]]; then
	NUM_CORES="2"
fi

mkdir -p $INSTALL_ROOT

OPENSSL_VERSION="1.0.2f"
NGINX_VERSION="1.8.1"
ERLANG_VERSION="18.2.1"
RABBIT_VERSION="3.6.0"
MONO_VERSION="4.2.2"
MONO_MINOR_RELEASE="29"
PSQL_VERSION="9.5.0"
GO_VERSION="1.5.3"
MONGODB_VERSION="3.2.4"

./openssl_install.sh $INSTALL_ROOT $OPENSSL_VERSION $NUM_CORES
./nginx_install.sh $INSTALL_ROOT $NGINX_VERSION $INSTALL_ROOT/_openssl_ $NUM_CORES
./erlang_install.sh $INSTALL_ROOT $ERLANG_VERSION $NUM_CORES
./rabbitmq_install.sh $INSTALL_ROOT $RABBIT_VERSION
./mono_install.sh $INSTALL_ROOT $MONO_VERSION $MONO_MINOR_RELEASE $NUM_CORES
./postgresql_install.sh $INSTALL_ROOT $PSQL_VERSION $NUM_CORES
./go_install.sh $INSTALL_ROOT $GO_VERSION
./mongo_install.sh $INSTALL_ROOT $MONGODB_VERSION

../Configuration/rabbitmq_setup.sh  $INSTALL_ROOT


export PATH=$PATH:$INSTALL_ROOT