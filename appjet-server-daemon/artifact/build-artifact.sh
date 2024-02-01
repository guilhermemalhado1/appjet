#!/bin/bash

# This script builds the appjet-client-cli executable

# Get today's date and time in the format MMDDYYYYHHMM
datetime=$(date +"%m%d%Y%H%M")

cd ../

# Create a build folder with the appended date and time
go build -o "../builds/appjet-server-daemon-$datetime/appjet-server-daemon" .

# create build inside for debug purposes - comment if needed below
go build -o "../appjet-server-daemon/artifact/builds/appjet-server-daemon-$datetime/appjet-server-daemon" .