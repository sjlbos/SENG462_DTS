#!/bin/bash

# Input parameters
USER=$1
PASSWORD=$2
BINARIES_ROOT=$3
REPO_ROOT=$4
ENVIRONMENT_ROOT=$5
DEPLOYMENT_ROOT=$6

if [[ -z "$USER" ]]; then
	echo "No username specified."
	exit 1
fi

if [[ -z "$PASSWORD" ]]; then
	echo "No password specified."
	exit 1
fi

if [[ -z "$BINARIES_ROOT" ]]; then
	echo "No binaries directory specified."
	exit 1
fi

if [[ -z "$REPO_ROOT" ]]; then
	echo "Repository root directory not specified."
	exit 1
fi

if [[ -z "$ENVIRONMENT_ROOT" ]]; then
	echo "Environment installtion directory not specified."
	exit 1
fi

if [[ -z "$DEPLOYMENT_ROOT" ]]; then
	echo "Deployment directory not specified."
	exit 1
fi

# Server Hosts

HOST_SUFFIX=".seng.uvic.ca"

DTS_MESSAGE_BROKER_SERVER=""
DTS_MESSAGE_BROKER_PORT=""

WLG_SLAVE_SERVERS=()
WLG_MESSAGE_BROKER_SERVER=""
WLG_MESSAGE_BROKER_PORT=""

WEB_SERVERS=()
WEB_SERVER_PORT=""

API_LOAD_BALANCER=""
API_LOAD_BALANCER_PORT=""

API_SERVERS=()
API_PORT=""

QUOTE_CACHE_SERVER=""
QUOTE_CACHE_PORT=""

TRANSACTION_MONITOR_SERVERS=()
TRANSACTION_MONITOR_PORT=""

TRIGGER_MANAGER_SERVER=""

DTS_DB_SERVER=""
DTS_DB_PORT=""

AUDIT_DB_SERVER=""
AUDIT_DB_PORT=""

# Host Environment Paths
DB_DATA_DIR=$ENVIRONMENT_ROOT/data


# Deploy DTS Database



# Deploy Audit Database



# Deploy Transaction Monitors
for host in "${TRANSACTION_MONITOR_SERVERS[@]}"
do

done

# Deploy Trigger Manager

# Deploy Quote Cache 

# Deploy APIs 
for host in "${API_SERVERS[@]}"
do

done

# Deploy Web Servers
for host in "${WEB_SERVERS[@]}"
do

done

# Deploy API Load Balancer



# Deploy WLG Slaves 