package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	waitRetries = 1000
	waitFactor  = 2
	waitDefault = 5
	waitRand    = 10
	
	k8ipprefix  = "k8ip"
)

func backoffExp(f func() error, errPre string) error {
	waitFor := waitDefault + rand.Intn(waitRand)
	for i := 0; i < waitRetries; i++ {
		err := f()
		if err != nil {
			fmt.Println(errPre, err.Error())
		} else {
			return nil
		}
		waitFor = waitFor * waitFactor
		fmt.Printf("Wait for: %d seconds and retry ...\n", waitFor)
		time.Sleep(time.Duration(waitFor) * time.Second)
	}
	return fmt.Errorf("Timeout reached")
}

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
	var vm compute.VirtualMachine
	var nic network.Interface
	var err error

	err = backoffExp(func() error {
		vm, err = vmClient.Get(groupName, vmName, compute.InstanceView)
		return err
	}, "Error getting VM details: ")
	if err != nil {
		return nil, err
	}
	fmt.Println("Obtained VM ID: ", *vm.ID)

	nicID := *(*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0].ID
	fmt.Println("NIC ID: ", nicID)

	nicName := path.Base(nicID)

	err = backoffExp(func() error {
		nic, err = nicClient.Get(groupName, nicName, "")
		return err
	}, "Error getting NIC details: ")
	if err != nil {
		return nil, err
	}
	return &nic, nil
}


func addIPstoVMNic(nicClient network.InterfacesClient, nic network.Interface, groupName string, count int) (err error) {
	// if run in parallel targetting the same VM/NIC, only one of the invocations will succeed and the others will fail.
	newidx := 0
	existingIPs := 0
	var primarySubnet network.Subnet

	for _, ipconfig := range *nic.InterfacePropertiesFormat.IPConfigurations {
		name := *ipconfig.Name
		fmt.Println("ipconfig name: ", name)

		if strings.HasPrefix(name, k8ipprefix) {
			if idx, err := strconv.Atoi(strings.TrimPrefix(name, k8ipprefix)); err == nil {
				existingIPs = existingIPs + 1
				if idx > newidx {
					fmt.Println("setting new index to: ", idx)
					newidx = idx
				}
			}
		}

		if ipconfig.InterfaceIPConfigurationPropertiesFormat == nil {
			continue
		}

		if (*ipconfig.InterfaceIPConfigurationPropertiesFormat).Primary == nil {
			continue
		}

		if *(*ipconfig.InterfaceIPConfigurationPropertiesFormat).Primary {
			primarySubnet = *(*ipconfig.InterfaceIPConfigurationPropertiesFormat).Subnet
		}
	}

	if existingIPs >= count {
		fmt.Printf("VM already has %d IPs. Skipping addition of new IPs\n", existingIPs)
		return nil
	}

	count = count - existingIPs
	for i := 0; i < count; i++ {
		newidx = newidx + 1
		newIPCfgName := fmt.Sprintf("%s%d", k8ipprefix, newidx)
		fmt.Println("New ipconfig name: ", newIPCfgName)
		newIPCfg := network.InterfaceIPConfiguration{
			Name: &newIPCfgName,
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &primarySubnet,
			},
		}

		*nic.InterfacePropertiesFormat.IPConfigurations = append(*nic.InterfacePropertiesFormat.IPConfigurations, newIPCfg)
	}

	fmt.Println("Waiting to update NIC ....")
	err = backoffExp(func() error {
		_, err = nicClient.CreateOrUpdate(groupName, *nic.Name, nic, nil)
		return err
	}, "Failed to update NIC: ")
	if err != nil {
		return nil
	}

	return nil
}
