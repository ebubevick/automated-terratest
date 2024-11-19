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

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	uniquePostfix := random.UniqueId()

	// Get the environment variables
	subscription := os.Getenv("subscription")
	clientID := os.Getenv("client_id")
	clientSecret := os.Getenv("client_secret")
	tenantID := os.Getenv("tenant_id")

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
	}

	// Cleanup after test
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and sanitize them
	resourceGroupName := strings.TrimSpace(terraform.Output(t, terraformOptions, "resource_group_name"))
	storageAccountName := strings.TrimSpace(terraform.Output(t, terraformOptions, "storage_account_name"))
	storageAccountTier := strings.TrimSpace(terraform.Output(t, terraformOptions, "storage_account_account_tier"))
	storageAccountKind := strings.TrimSpace(terraform.Output(t, terraformOptions, "storage_account_account_kind"))
	storageBlobContainerName := strings.TrimSpace(terraform.Output(t, terraformOptions, "storage_container_name"))
	storageFileShareName := strings.TrimSpace(terraform.Output(t, terraformOptions, "storage_fileshare_name"))

	// Verify storage account properties and ensure it matches the output.
	storageAccountExists := azure.StorageAccountExists(t, storageAccountName, resourceGroupName, subscription)
	assert.True(t, storageAccountExists, "storage account does not exist")

	containerExists := azure.StorageBlobContainerExists(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscription)
	assert.True(t, containerExists, "storage container does not exist")

	fileShareExists := azure.StorageFileShareExists(t, storageFileShareName, storageAccountName, resourceGroupName, "")
	assert.True(t, fileShareExists, "File share does not exist")

	publicAccess := azure.GetStorageBlobContainerPublicAccess(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscription)
	assert.False(t, publicAccess, "storage container has public access")

	accountKind := azure.GetStorageAccountKind(t, storageAccountName, resourceGroupName, subscription)
	assert.Equal(t, storageAccountKind, accountKind, "storage account kind mismatch")

	skuTier := azure.GetStorageAccountSkuTier(t, storageAccountName, resourceGroupName, subscription)
	assert.Equal(t, storageAccountTier, skuTier, "sku tier mismatch")

	actualDNSString := azure.GetStorageDNSString(t, storageAccountName, resourceGroupName, subscription)
	storageSuffix, _ := azure.GetStorageURISuffixE()
	expectedDNS := fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix)
	assert.Equal(t, expectedDNS, actualDNSString, "Storage DNS string mismatch")
}
