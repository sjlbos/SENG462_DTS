#!/bin/bash

# This deploys a DTS database on a remote server.
# Execute using:
#
#     ./deploy_database.sh 
#			{your username} 
#			{server name} 
#			{database name}
#			{database port}
#			{database data directory on remote server}
#			{directory containing database creation script}
#			{database creation script name}

USER=$1
SERVER=$2
DB_NAME=$3
DB_PORT=$4
DATA_DIR=$5
CREATE_SCRIPT_DIR=$6
CREATE_SCRIPT_NAME=$7

if [[ -z "$USER" ]]; then
	echo "No user specified."
	exit 1
fi

if [[ -z "$SERVER" ]]; then
	echo "No server specified."
	exit 1
fi

if [[ -z "$DB_NAME" ]]; then
	echo "No database name specified."
	exit 1
fi

if [[ -z "$DB_PORT" ]]; then
	echo "No database port specified."
	exit 1
fi

if [[ -z "$DATA_DIR" ]]; then
	echo "No data directory specified."
	exit 1
fi

if [[ -z "$CREATE_SCRIPT_DIR" ]]; then
	echo "Database creation script directory not specified."
	exit 1
fi

if [[ -z "$CREATE_SCRIPT_NAME" ]]; then
	echo "Database creation script name not specified."
	exit 1
fi

SSH_PATH=$USER@$SERVER
CREATE_SCRIPT_REMOTE_DIR=dts_tmp
CREATE_SCRIPT_REMOTE_PATH=$CREATE_SCRIPT_REMOTE_DIR/$CREATE_SCRIPT_NAME

# Create temp directory on server
ssh $SSH_PATH "mkdir -p $CREATE_SCRIPT_REMOTE_DIR"

# Copy database create script to server
scp $CREATE_SCRIPT_DIR/$CREATE_SCRIPT_NAME $SSH_PATH:$CREATE_SCRIPT_REMOTE_DIR

# Peform database deployment
ssh $SSH_PATH <<EOF
	dropdb -p $DB_PORT $DB_NAME
	pg_ctl stop -w -D $DATA_DIR
	pg_ctl start -w -D $DATA_DIR -l $DATA_DIR/logfile.txt
	createdb -p $DB_PORT $DB_NAME
	psql -d $DB_NAME -p $DB_PORT -U dts_user -f $CREATE_SCRIPT_REMOTE_PATH
	rm -rf CREATE_SCRIPT_REMOTE_DIR
EOF
