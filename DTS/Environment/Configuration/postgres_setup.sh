#!/bin/bash

# This script configures and starts a new installation of a Postgresql server
# Execute using:
#
#     postgres_setup.sh {install directory} {data directory} {port} {dts user password}

INSTALL_DIR=$1
DATA_DIR=$2
PORT=$3
DTS_USER_PASSWORD=$4

if [[ -z "$INSTALL_DIR" ]]; then
	echo "No install directory specified."
	exit 1
fi

if [[ -z "$DATA_DIR" ]]; then
	echo "No data directory specified."
	exit 1
fi

if [[ -z "$PORT" ]]; then
	echo "No database port specified."
	exit 1
fi

if [[ -z "$DTS_USER_PASSWORD" ]]; then
	echo "No password specified for dts_user."
	exit 1
fi

# Initialize data directory
mkdir -p $DATA_DIR
$INSTALL_DIR/bin/initdb -D $DATA_DIR

# Allow password authentication
echo "local samerole all md5" >> $DATA_DIR/pg_hba.conf
echo "host samerole all samenet md5" >> $DATA_DIR/pg_hba.conf

# Configure port
sed -i "s/#port = 5432/port = $PORT/" $DATA_DIR/postgresql.conf

# Enable remote access
sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" $DATA_DIR/postgresql.conf

# Start server
$INSTALL_DIR/bin/pg_ctl start -w -D $DATA_DIR -l $DATA_DIR/logfile.txt

# Create dts_user account
echo "CREATE USER dts_user WITH PASSWORD '$DTS_USER_PASSWORD';" | $INSTALL_DIR/bin/psql -d postgres -p $PORT
echo "ALTER USER dts_user CREATEDB;" | $INSTALL_DIR/bin/psql -d postgres -p $PORT
