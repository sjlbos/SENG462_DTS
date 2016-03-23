#!/bin/bash

INSTALL_ROOT=$1
DATA_DIR=$2

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "Mongo installation directory not set."
	exit 1
fi

if [[ -z "$DATA_DIR" ]]; then
	echo "Data directory not set."
	exit 1
fi

mkdir -p $DATA_DIR
mkdir -p $DATA_DIR/db
mkdir -p $DATA_DIR/log

cat << EOF > $DATA_DIR/mongod.cfg
systemLog: 
    destination: file 
    path: ./test/data/log/mongod.log 
storage: 
    dbPath: ./test/data/db 
net: 
    port: 44410 
    bindIp: 127.0.0.1 
    maxIncomingConnections: 65536 
EOF

$INSTALL_ROOT/bin/mongod --fork --config $DATA_DIR/mongod.cfg

