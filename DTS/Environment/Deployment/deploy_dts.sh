#!/bin/bash

# Input parameters
USER=$1
REPO_ROOT=$2
ENVIRONMENT_ROOT=$3
DEPLOYMENT_ROOT=$4

if [[ -z "$USER" ]]; then
	echo "No username specified."
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

##########################################################################################

# Server Hosts

HOST_SUFFIX=".seng.uvic.ca"

WEB_SERVERS=("")
WEB_SERVER_PORT=""

API_SERVERS=("b148" "b149" "b150")
API_PORT="44410"

QUOTE_CACHE_SERVER="b143"
QUOTE_CACHE_PORT="44410"

TRANSACTION_MONITOR_SERVERS=("b136")
TRANSACTION_MONITOR_PORT="44410"

TRIGGER_MANAGER_SERVER="b135"

DTS_DB_SERVERS=("b133" "b142" "b144")
DTS_DB_PORT="44410"

AUDIT_DB_SERVER="b132"
AUDIT_DB_PORT="44410"

WLG_SLAVE_SERVERS=("b138" "b139" "b140" "b141" "b145" "b146" "b147")

# Local Paths
PACKAGE_DIR=$REPO_ROOT/packages
BUILD_DIR=$REPO_ROOT/bin

# Remote Paths
DB_DATA_DIR=$ENVIRONMENT_ROOT/data

##########################################################################################

function deployZipFile {
	SSH_PATH=$1
	DEPLOYMENT_DIR=$2
	FILE_DIR=$3
	FILE_NAME=$4

	ssh $SSH_PATH "mkdir -p $DEPLOYMENT_DIR"
	scp $FILE_DIR/$FILE_NAME $SSH_PATH:$DEPLOYMENT_DIR
	ssh $SSH_PATH <<EOF
	unzip -o $DEPLOYMENT_DIR/$FILE_NAME -d $DEPLOYMENT_DIR
	rm -f $DEPLOYMENT_DIR/$FILE_NAME
	chmod 770 -R $DEPLOYMENT_DIR/${FILE_NAME%.*}
EOF
}

##########################################################################################

# Deploy DTS Databases
for host in "${DTS_DB_SERVERS[@]}" 
do
	echo "Deploying DTS Database to $host."
	$REPO_ROOT/DTS/Environment/Deployment/deploy_database.sh $USER $host$HOST_SUFFIX dts $DTS_DB_PORT $DB_DATA_DIR $REPO_ROOT/DTS/Database CreateDtsDb.sql
done

# Deploy Audit Database
echo "Deploying DTS Audit Database"
$REPO_ROOT/DTS/Environment/Deployment/deploy_database.sh $USER $AUDIT_DB_SERVER$HOST_SUFFIX dts_audit $AUDIT_DB_PORT $DB_DATA_DIR $REPO_ROOT/DTS/Database CreateDtsAuditDb.sql

# Build DTS
echo "Building DTS binaries..."
$REPO_ROOT/DTS/Build/build.sh $REPO_ROOT
echo "Build complete."

# Package Binaries
echo "Creating DTS packages..."

rm -rf $PACKAGE_DIR
mkdir $PACKAGE_DIR

cd $BUILD_DIR
zip -r $PACKAGE_DIR/TransactionMonitor.zip TransactionMonitor
zip -r $PACKAGE_DIR/TriggerManager.zip TriggerManager
zip -r $PACKAGE_DIR/WorkloadGeneratorSlave.zip WorkloadGeneratorSlave
zip -r $PACKAGE_DIR/DtsApi.zip DtsApi
zip -r $PACKAGE_DIR/QuoteCache.zip QuoteCache
zip -r $PACKAGE_DIR/QuoteRunner.zip QuoteRunner
echo "Package creation complete."

# Deploy Transaction Monitors
for host in "${TRANSACTION_MONITOR_SERVERS[@]}"
do
	echo "Deploying Transaction Monitor to $host."
	deployZipFile "$USER@$host$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR TransactionMonitor.zip
done

# Deploy Trigger Manager
echo "Deploying Trigger Manager to $TRIGGER_MANAGER_SERVER."
deployZipFile "$USER@$TRIGGER_MANAGER_SERVER$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR TriggerManager.zip

# Deploy Quote Cache
echo "Deploying Quote Cache to $QUOTE_CACHE_SERVER." 
deployZipFile "$USER@$QUOTE_CACHE_SERVER$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR QuoteCache.zip

# Deploy Quote Runners
for host in "${QUOTE_RUNNER_SERVERS[@]}"
do
	echo "Deploying Quote Runner to $host."
	deployZipFile "$USER@$host$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR QuoteRunner.zip
done

# Deploy APIs 
for host in "${API_SERVERS[@]}"
do
	echo "Deploying DTS API to $host."
	deployZipFile "$USER@$host$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR DtsApi.zip
done

# Deploy Web Servers
for host in "${WEB_SERVERS[@]}"
do
	echo "Deploying web server to $host."
	#deployZipFile "$USER@$host$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR {Package Name Here!}
done

# Deploy WLG Slaves
for host in "${WLG_SLAVE_SERVERS[@]}"
do
	echo "Deploying WLG slave to $host."
	deployZipFile "$USER@$host$HOST_SUFFIX" $DEPLOYMENT_ROOT $PACKAGE_DIR WorkloadGeneratorSlave.zip
done
