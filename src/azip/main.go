package main

import (
	"fmt"
	"os"
)

func validate() {
	if !checkEnvVars("APP_ID", "APP_SECRET", "TENANT_ID", "ACCOUNT_ID", "GROUP_NAME", "VM_NAME") {
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

func main() {
	validate()
	env := map[string]string{
		"AZURE_CLIENT_ID":       os.Getenv("APP_ID"),
		"AZURE_CLIENT_SECRET":   os.Getenv("APP_SECRET"),
		"AZURE_SUBSCRIPTION_ID": os.Getenv("ACCOUNT_ID"),
		"AZURE_TENANT_ID":       os.Getenv("TENANT_ID"),
		"AZURE_GROUP_NAME":      os.Getenv("GROUP_NAME"),
		"AZURE_VM_NAME":		 os.Getenv("VM_NAME"),
	}
	nicClient, vmClient := initClients(env)
	
	nic, err := getVMNic(vmClient, nicClient, env["AZURE_GROUP_NAME"], env["AZURE_VM_NAME"])
	if err != nil {
		os.Exit(1)
	}
	err = addIPtoVMNic(nicClient, *nic, env["AZURE_GROUP_NAME"])
}
