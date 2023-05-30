#!/bin/bash

# This is used to allow to run and test locally 
if [ -z $BIN_FILE ]; then 
    BIN_FILE=/opt/scream/bin/scream 
fi 

echo "starting service ${BIN_FILE}"

$BIN_FILE stop
$BIN_FILE daemon


if [ $? -ne 0 ]; then 
    echo "Failed to start service ${BIN_FILE}"
    exit 1
fi