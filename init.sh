#!/bin/bash

export VM_NAME=$(curl -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq .compute.name | sed -e 's/^"//' -e 's/"$//')
echo "vm_name: $vm_name"

export GROUP_NAME=$(curl -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq .compute.resourceGroupName | sed -e 's/^"//' -e 's/"$//')
echo "rg_name: $rg_name"

export IP_COUNT=20

echo "invoking azip with $VM_NAME $GROUP_NAME $IP_COUNT."

azip
