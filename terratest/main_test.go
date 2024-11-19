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

	// Generate a unique postfix for resource names
	uniquePostfix := random.UniqueId()

	// Get environment variables and validate they are set
	subscription := os.Getenv("ARM_SUBSCRIPTION_ID")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantID := os.Getenv("ARM_TENANT_ID")

	if subscription == "" || clientID == "" || clientSecret == "" || tenantID == "" {
		t.Fatal("Required environment variables (ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID) are not set")
	}

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

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)

	// Retrieve outputs from Terraform
	resourceGroupName := getTerraformOutput(t, terraformOptions, "resource_group_name")
	storageAccountName := getTerraformOutput(t, terraformOptions, "storage_account_name")
	storageAccountTier := getTerraformOutput(t, terraformOptions, "storage_account_account_tier")
	storageAccountKind := getTerraformOutput(t, terraformOptions, "storage_account_account_kind")
	storageBlobContainerName := getTerraformOutput(t, terraformOptions, "storage_container_name")
	storageFileShareName := getTerraformOutput(t, terraformOptions, "storage_fileshare_name")

	// Verify storage account properties
	assert.True(t, azure.StorageAccountExists(t, storageAccountName, resourceGroupName, subscription), "Storage account does not exist")
	assert.True(t, azure.StorageBlobContainerExists(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscription), "Storage container does not exist")
	assert.True(t, azure.StorageFileShareExists(t, storageFileShareName, storageAccountName, resourceGroupName, subscription), "File share does not exist")

	// Verify storage account DNS
	storageSuffix, _ := azure.GetStorageURISuffixE()
	expectedDNS := fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix)
	actualDNSString := azure.GetStorageDNSString(t, storageAccountName, resourceGroupName, subscription)
	assert.Equal(t, expectedDNS, actualDNSString, "Storage DNS string mismatch")

	// Verify storage account kind and tier
	assert.Equal(t, storageAccountKind, azure.GetStorageAccountKind(t, storageAccountName, resourceGroupName, subscription), "Storage account kind mismatch")
	assert.Equal(t, storageAccountTier, azure.GetStorageAccountSkuTier(t, storageAccountName, resourceGroupName, subscription), "Storage account tier mismatch")

	// Verify public access
	assert.False(t, azure.GetStorageBlobContainerPublicAccess(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscription), "Storage container has public access")
}

// Helper function to safely retrieve Terraform output
func getTerraformOutput(t *testing.T, options *terraform.Options, key string) string {
	output, err := terraform.OutputE(t, options, key)
	if err != nil {
		t.Fatalf("Failed to get Terraform output for key '%s': %s", key, err)
	}
	return output
}
