#!/bin/bash

set -e

if [ "$1" == "cleanup" ]; then
    echo "Cleaning up resources..."
    docker-compose -f ${PWD}/gateway-application/docker-compose.yml down
    ${PWD}/test-network/network.sh down

elif [ "$1" == "restart" ]; then
    echo "Restarting Application.."
    docker-compose -f ${PWD}/gateway-application/docker-compose.yml down
    docker-compose -f ${PWD}/gateway-application/docker-compose.yml up --build    
else
    echo "Setting up the hyperledger and the application..."
    ${PWD}/test-network/network.sh up createChannel -c mychannel -ca
    ${PWD}/test-network/network.sh deployCC -ccn basic -ccp ../chaincode-bank -ccl go
    docker-compose -f ${PWD}/gateway-application/docker-compose.yml up --build
fi
