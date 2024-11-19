package test

import (
	"encoding/json"
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

	// Validate environment variables
	envVars := []string{"subscription", "client_id", "client_secret", "tenant_id"}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Environment variable %s must be set", envVar)
		}
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

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and sanitize them
	resourceGroupName := sanitizeTerraformOutput(t, terraform.Output(t, terraformOptions, "resource_group_name"))
	storageAccountName := sanitizeTerraformOutput(t, terraform.Output(t, terraformOptions, "storage_account_name"))
	storageAccountTier := sanitizeTerraformOutput(t, terraform.Output(t, terraformOptions, "storage_account_account_tier"))
	storageAccountKind := sanitizeTerraformOutput(t, terraform.Output(t, terraformOptions, "storage_account_account_kind"))
	storageBlobContainerName := sanitizeTerraformOutput(t, terraform.Output(t, terraformOptions, "storage_container_name"))
	storageFileShareName := sanitizeTerraformOutput(t, terraform.Output(t, terraformOptions, "storage_fileshare_name"))

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

// sanitizeTerraformOutput trims and validates the JSON output.
func sanitizeTerraformOutput(t *testing.T, output string) string {
	output = strings.TrimSpace(output)
	var sanitizedOutput interface{}
	if err := json.Unmarshal([]byte(output), &sanitizedOutput); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}
	sanitizedBytes, err := json.Marshal(sanitizedOutput)
	if err != nil {
		t.Fatalf("Failed to re-marshal JSON output: %v", err)
	}
	return string(sanitizedBytes)
}
