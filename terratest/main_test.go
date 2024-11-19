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
	// subscriptionID := ""
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

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	storageAccountName := terraform.Output(t, terraformOptions, "storage_account_name")
	storageAccountTier := terraform.Output(t, terraformOptions, "storage_account_account_tier")
	storageAccountKind := terraform.Output(t, terraformOptions, "storage_account_account_kind")
	storageBlobContainerName := terraform.Output(t, terraformOptions, "storage_container_name")
	storageFileShareName := terraform.Output(t, terraformOptions, "storage_fileshare_name")

	// website::tag::4:: Verify storage account properties and ensure it matches the output.
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
