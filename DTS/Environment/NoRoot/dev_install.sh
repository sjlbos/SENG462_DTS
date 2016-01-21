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

ERLANG_VERSION="18.2.1"
RABBIT_VERSION="3.6.0"
MONO_VERSION="4.2.2"
MONO_MINOR_RELEASE="29"
PSQL_VERSION="9.5.0"

source ./erlang_install.sh $INSTALL_ROOT $ERLANG_VERSION $NUM_CORES
source ./rabbitmq_install.sh $INSTALL_ROOT $RABBIT_VERSION
source ./mono_install.sh $INSTALL_ROOT $MONO_VERSION $MONO_MINOR_RELEASE $NUM_CORES
source ./postgresql_install.sh $INSTALL_ROOT $PSQL_VERSION $NUM_CORES


