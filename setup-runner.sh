#!/bin/bash

if [ $# -lt 1 ]; then 
    echo "Usage: ./setup-runner [token]"
    exit 1
fi

token=$1

# Create a folder
mkdir actions-runner && cd actions-runner

# Download the latest runner package
curl -o actions-runner-linux-x64-2.334.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.334.0/actions-runner-linux-x64-2.334.0.tar.gz

# Extract the installer
tar xzf ./actions-runner-linux-x64-2.334.0.tar.gz

# Create the runner and start the configuration experience
./config.sh --url https://github.com/fakeowl1/lab1-devops --token $token --unattended
