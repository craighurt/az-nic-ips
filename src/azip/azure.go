package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/go-autorest/autorest/azure"
)

func initClients(env map[string]string) (network.InterfacesClient, compute.VirtualMachinesClient) {
	rmEndpoint := azure.PublicCloud.ResourceManagerEndpoint
	// handle other endpoints like Azure Gov/China/etc
	if uri := os.Getenv("RESOURCE_MANAGER_ENDPOINT"); uri != "" {
		rmEndpoint = uri
	}

	spt, err := helpers.NewServicePrincipalTokenFromCredentials(env, rmEndpoint)
	if err != nil {
		fmt.Printf("ERROR: Getting SP token - check that all ENV variables are set")
		os.Exit(1)
	}

	// Create Network Interface Client
	nicClient := network.NewInterfacesClientWithBaseURI(rmEndpoint, env["AZURE_SUBSCRIPTION_ID"])
	nicClient.Authorizer = spt
	// Create VM Client
	vmClient := compute.NewVirtualMachinesClientWithBaseURI(rmEndpoint, env["AZURE_SUBSCRIPTION_ID"])
	vmClient.Authorizer = spt

	return nicClient, vmClient
}

func getVMNic(vmClient compute.VirtualMachinesClient, nicClient network.InterfacesClient, groupName, vmName string) (*network.Interface, error) {
	vm, err := vmClient.Get(groupName, vmName, compute.InstanceView)
	if err != nil {
		fmt.Println("Error getting VM details: ", err.Error())
		return nil, err
	}
	fmt.Println("Obtained VM ID: ", *vm.ID)
	
	nicID := *(*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0].ID
	fmt.Println("NIC ID: ", nicID)
	
	nicName := path.Base(nicID)
	
	nic, err := nicClient.Get(groupName, nicName, "")
	if err != nil {
		fmt.Println("Error getting NIC details: ", err.Error())
		return nil, err
	}
	return &nic, nil
}

func addIPtoVMNic(nicClient network.InterfacesClient, nic network.Interface, groupName string) (err error) {
	newidx := 0
	var primarySubnet network.Subnet
	
	for _, ipconfig := range *nic.InterfacePropertiesFormat.IPConfigurations {
		name := *ipconfig.Name
		fmt.Println("ipconfig name: ", name)
		
		if (strings.HasPrefix(name, "ipconfig")) {
			if idx, err := strconv.Atoi(strings.TrimPrefix(name, "ipconfig")); err == nil {
				if (idx > newidx) {
					fmt.Println("setting new index to: ", idx)
					newidx = idx
				}
			}
		}
		
		if (*(*ipconfig.InterfaceIPConfigurationPropertiesFormat).Primary) {
			primarySubnet = *(*ipconfig.InterfaceIPConfigurationPropertiesFormat).Subnet
		}
	}
	newidx = newidx + 1
	newIPCfgName := fmt.Sprintf("ipconfig%d", newidx)
	fmt.Println("New ipcfg name: ", newIPCfgName)
	
	newIPCfg := network.InterfaceIPConfiguration {
		Name: &newIPCfgName,
		InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: network.Dynamic,
			Subnet: &primarySubnet,
		},
	}
	
	*nic.InterfacePropertiesFormat.IPConfigurations = append(*nic.InterfacePropertiesFormat.IPConfigurations, newIPCfg)
	
	_, err = nicClient.CreateOrUpdate(groupName, *nic.Name, nic, nil)
	fmt.Println("Waiting to update NIC ....") 
	if err != nil {
		fmt.Println("Failed to update NIC: ", err.Error())
		return err
	}
	return nil
}
