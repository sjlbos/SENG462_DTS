#!/bin/bash

INSTALL_ROOT=$1
DATA_DIR=$2
PORT=$3

if [[ -z "$INSTALL_ROOT" ]]; then
	echo "Mongo installation directory not set."
	exit 1
fi

if [[ -z "$DATA_DIR" ]]; then
	echo "Data directory not set."
	exit 1
fi

if [[ -z "$PORT" ]]; then
    echo "MongoDb port not specified."
    exit 1
fi

mkdir -p $DATA_DIR
mkdir -p $DATA_DIR/db
mkdir -p $DATA_DIR/log

cat << EOF > $DATA_DIR/mongod.cfg
systemLog: 
    destination: file 
    path: $DATA_DIR/log/mongod.log 
storage: 
    dbPath: $DATA_DIR/db 
net: 
    port: $PORT
    bindIp: 127.0.0.1 
    maxIncomingConnections: 500
EOF

$INSTALL_ROOT/bin/mongod --fork --config $DATA_DIR/mongod.cfg

