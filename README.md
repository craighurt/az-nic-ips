
[![GoDoc](https://godoc.org/github.com/ddebroy/az-nic-ips?status.svg)](https://godoc.org/github.com/ddebroy/az-nic-ips)
[![TravisCI](https://travis-ci.org/ddebroy/az-nic-ips.svg?branch=master)](https://travis-ci.org/ddebroy/az-nic-ips)
[![Go Report Card](https://goreportcard.com/badge/github.com/ddebroy/az-nic-ips)](https://goreportcard.com/report/github.com/ddebroy/az-nic-ips)
[![Badge Badge](http://doyouevenbadge.com/github.com/ddebroy/az-nic-ips)](http://doyouevenbadge.com)

# AZ NIC IPs
tool to add IPs to NICs in Azure VMs

To build the binary, simply run `make`.

To wrap up the binary in a docker alpine container, run `make image`


## Deploy as a single container
To run the container, launch it with:
```
$ docker run -e APP_ID={AZURE_SP_ID} -e APP_SECRET={AZURE_SP_SECRET} -e ACCOUNT_ID={AZURE_SUBSCRIPTION_ID} -e TENANT_ID={AZURE_TENANT_ID} -e GROUP_NAME={RESOURCE_GROUP_NAME_WITH_VM} -e VM_NAME={NAME_OF_VM_IN_RG} docker4x/az-nic-ips:latest
```

You should see a new IP added to the first NIC of VM with VM_NAME

## Deploy as a Docker service

1. Create a `azure_ucp_admin.toml` file with the following information:
```
AZURE_CLIENT_ID = "<AD App ID >"
AZURE_TENANT_ID = "<AD Tenant ID >"
AZURE_SUBSCRIPTION_ID = "<Azure subscription ID>"
AZURE_CLIENT_SECRET = "<AD App Secret >"
```

2. Create a Docker Swarm secret based on the "azure_ucp_admin.toml" file
```
docker secret create azure_ucp_admin.toml azure_ucp_admin.toml
```

3. Create a global swarm service using this image
```
docker service create \
  --mode=global \
  --secret=azure_ucp_admin.toml \
  --log-driver json-file \
  --log-opt max-size=1m \
  --env IPCOUNT=128 \
  --name ipallocator \
  docker4x/az-nic-ips:latest
```
