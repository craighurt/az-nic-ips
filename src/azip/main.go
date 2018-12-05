package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	AzureClientID       string `json:"aadClientId"`
	AzureTenantID       string `json:"tenantId"`
	AzureSubscriptionID string `json:"subscriptionId"`
	AzureClientSecret   string `json:"aadClientSecret"`
	AzureCloudName      string `json:"cloud"`
}

func main() {
	var config AzureConfig
	file, err := ioutil.ReadFile("/run/secrets/azure_ucp_admin.json")
	if err != nil {
		fmt.Printf("ERROR: failed to read file: %v\n", err)
		return
	}
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("ERROR: could not decode secrets file %v", err)
		return
	}

	env := map[string]string{
		"AZURE_CLOUD_NAME":      config.AzureCloudName,
		"AZURE_CLIENT_ID":       config.AzureClientID,
		"AZURE_CLIENT_SECRET":   config.AzureClientSecret,
		"AZURE_SUBSCRIPTION_ID": config.AzureSubscriptionID,
		"AZURE_TENANT_ID":       config.AzureTenantID,
		"AZURE_GROUP_NAME":      os.Getenv("GROUP_NAME"),
		"AZURE_VM_NAME":         os.Getenv("VM_NAME"),
		"IP_COUNT":              os.Getenv("IP_COUNT"),
	}
	nicClient, vmClient := initClients(env)

	vm, err := getVM(vmClient, env["AZURE_VM_NAME"], env["AZURE_GROUP_NAME"])
	if vm == nil || err != nil {
		os.Exit(1)
	}

	if skipVM(*vm) {
		fmt.Println("Skipping VM")
		os.Exit(0)
	}

	nic, err := getNIC(nicClient, *vm, env["AZURE_GROUP_NAME"])
	if nic == nil || err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ips, err := strconv.Atoi(env["IP_COUNT"])
	if err != nil {
		fmt.Println("ERROR: Invalid IP_COUNT specified")
		os.Exit(1)
	}

	err = addIPstoVMNic(nicClient, *nic, env["AZURE_GROUP_NAME"], ips)
	if err != nil {
		fmt.Println("ERROR: failed to add IPs to VM")
		os.Exit(1)
	}
}
