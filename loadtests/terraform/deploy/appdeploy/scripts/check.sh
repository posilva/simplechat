#!/bin/bash 
set -ex -o pipefail

# This is used to allow to run and test locally 
if [ -z $BIN_FILE ]; then 
    BIN_FILE=/opt/scream/bin/scream 
fi 

$BIN_FILE pid 

if [ $? -ne 0 ]; then 
    echo "Failed to start service ${BIN_FILE}"
    exit 1
fi