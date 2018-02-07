#!/bin/bash

while true
do
    export VM_NAME=$(curl -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq .compute.name | sed -e 's/^"//' -e 's/"$//')
    export GROUP_NAME=$(curl -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq .compute.resourceGroupName | sed -e 's/^"//' -e 's/"$//')

    if [[ -z "${IP_COUNT}" ]]; then
        export IP_COUNT=20
    fi

    echo "invoking azip with $VM_NAME $GROUP_NAME $IP_COUNT."
    azip
    sleep 10
done
