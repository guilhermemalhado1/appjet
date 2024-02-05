#!/bin/bash

# This script builds the appjet-client-cli executable

# Get today's date and time in the format MMDDYYYYHHMM
datetime=$(date +"%m%d%Y%H%M")

cd ../

# create build inside for debug purposes - comment if needed below
go build -o "../appjet-cli/artifact/builds/appjet-cli-$datetime/appjet-cli" .