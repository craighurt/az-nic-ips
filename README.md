# azip
tool to add IPs to NICs in Azure VMs

To build the binary, simply run `make`.

To wrap up the binary in a docker alpine container, run `docker build -t azip -f Dockerfile .`

To run the container, launch it with `docker run -e APP_ID={AZURE_SP_ID} -e APP_SECRET={AZURE_SP_SECRET} -e ACCOUNT_ID={AZURE_SUBSCRIPTION_ID} -e TENANT_ID={AZURE_TENANT_ID} -e GROUP_NAME={RESOURCE_GROUP_NAME_WITH_VM} -e VM_NAME={NAME_OF_VM_IN_RG} azip`
