package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureStorage(t *testing.T) {
	t.Parallel()

	// Use random postfix to ensure unique resources
	uniquePostfix := random.UniqueId()

	// Get the environment variables
	subscription := os.Getenv("subscription")
	clientID := os.Getenv("client_id")
	clientSecret := os.Getenv("client_secret")
	tenantID := os.Getenv("tenant_id")

	// Ensure the variables are set
	assert.NotEmpty(t, subscription)
	assert.NotEmpty(t, clientID)
	assert.NotEmpty(t, clientSecret)
	assert.NotEmpty(t, tenantID)

	// Configure Terraform
	terraformOptions := &terraform.Options{
		// Path to Terraform code
		TerraformDir: "../",

		// Variables to pass to Terraform code
		Vars: map[string]interface{}{
			"subscription":  subscription,
			"client_id":     clientID,
			"client_secret": clientSecret,
			"tenant_id":     tenantID,
			"postfix":       strings.ToLower(uniquePostfix),
		},

		// Disable prompt for input during tests
		NoColor: true,
	}

	// Cleanup after test
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)

	// Validate outputs
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	storageAccountName := terraform.Output(t, terraformOptions, "storage_account_name")
	storageAccountTier := terraform.Output(t, terraformOptions, "storage_account_account_tier")
	storageAccountKind := terraform.Output(t, terraformOptions, "storage_account_account_kind")
	storageBlobContainerName := terraform.Output(t, terraformOptions, "storage_container_name")
	storageFileShareName := terraform.Output(t, terraformOptions, "storage_fileshare_name")

	// Assertions using Azure SDK
	storageAccountExists := azure.StorageAccountExists(t, storageAccountName, resourceGroupName, "")
	assert.True(t, storageAccountExists, "storage account does not exist")

	containerExists := azure.StorageBlobContainerExists(t, storageBlobContainerName, storageAccountName, resourceGroupName, "")
	assert.True(t, containerExists, "storage container does not exist")

	fileShareExists := azure.StorageFileShareExists(t, storageFileShareName, storageAccountName, resourceGroupName, "")
	assert.True(t, fileShareExists, "File share does not exist")

	publicAccess := azure.GetStorageBlobContainerPublicAccess(t, storageBlobContainerName, storageAccountName, resourceGroupName, "")
	assert.False(t, publicAccess, "storage container has public access")

	accountKind := azure.GetStorageAccountKind(t, storageAccountName, resourceGroupName, "")
	assert.Equal(t, storageAccountKind, accountKind, "storage account kind mismatch")

	skuTier := azure.GetStorageAccountSkuTier(t, storageAccountName, resourceGroupName, "")
	assert.Equal(t, storageAccountTier, skuTier, "sku tier mismatch")

	actualDNSString := azure.GetStorageDNSString(t, storageAccountName, resourceGroupName, "")
	storageSuffix, _ := azure.GetStorageURISuffixE()
	expectedDNS := fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix)
	assert.Equal(t, expectedDNS, actualDNSString, "Storage DNS string mismatch")
}
