#!/bin/bash

# This script builds the appjet-client-cli executable

# Get today's date and time in the format MMDDYYYYHHMM
datetime=$(date +"%m%d%Y%H%M")

cd ../

# Create a build folder with the appended date and time
GOOS=linux go build -o "../builds/appjet-decision-manager-$datetime/appjet-decision-manager" .

# create build inside for debug purposes - comment if needed below
GOOS=linux go build -o "../appjet-decision-manager/artifact/builds/appjet-decision-manager-$datetime/appjet-decision-manager" .