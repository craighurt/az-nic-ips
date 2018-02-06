package main

import (
	"fmt"
	"os"
	"strconv"
)

func validate() {
	if !checkEnvVars("GROUP_NAME", "VM_NAME", "IP_COUNT") {
		os.Exit(1)
	}
}

func checkEnvVars(envVars ...string) bool {
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("ERROR: Missing environment variable: %s.\n", envVar)
			return false
		}
	}
	return true
}

type AzureConfig struct {
	AzureClientID       string `toml:"AZURE_CLIENT_ID"`
	AzureTenantID       string `toml:"AZURE_TENANT_ID"`
	AzureSubscriptionID string `toml:"AZURE_SUBSCRIPTION_ID"`
	AzureClientSecret   string `toml:"AZURE_CLIENT_SECRET"`
}

func main() {
	/*
	var config AzureConfig
	if _, err := toml.DecodeFile("/run/secrets/azure_ucp_admin.toml", &config); err != nil {
		fmt.Printf("Error decoding file %v", err)
		return
	}
	*/

	env := map[string]string{
		"AZURE_CLIENT_ID":       os.Getenv("APP_ID"),
		"AZURE_CLIENT_SECRET":   os.Getenv("APP_SECRET"),
		"AZURE_SUBSCRIPTION_ID": os.Getenv("ACCOUNT_ID"),
		"AZURE_TENANT_ID":       os.Getenv("TENANT_ID"),
		"AZURE_GROUP_NAME":      os.Getenv("GROUP_NAME"),
		"AZURE_VM_NAME":         os.Getenv("VM_NAME"),
		"IP_COUNT":              os.Getenv("IP_COUNT"),
	}
	nicClient, _, vmssClient, vmssVMClient := initClients(env)

	nic, err := getVMSSNic(vmssClient, vmssVMClient, nicClient, env["AZURE_GROUP_NAME"], env["AZURE_VM_NAME"])
	if err != nil {
		os.Exit(1)
	}

	ips, err := strconv.Atoi(env["IP_COUNT"])
	if err != nil {
		fmt.Println("Invalid IP_COUNT specified")
		os.Exit(1)
	}
	if nic == nil {
		fmt.Printf("NIC for VM %s could not be located\n", env["AZURE_VM_NAME"])
		os.Exit(1)		
	}
	err = addIPstoVMNic(nicClient, *nic, env["AZURE_GROUP_NAME"], ips)

}
