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
	subscriptionID := "4bf3e463-ed9f-4148-8906-3eed094e0794"
	uniquePostfix := random.UniqueId()

	fmt.Println("AZURE_SUBSCRIPTION_ID:", os.Getenv("AZURE_SUBSCRIPTION_ID"))
	fmt.Println("AZURE_CLIENT_ID:", os.Getenv("AZURE_CLIENT_ID"))
	fmt.Println("AZURE_CLIENT_SECRET:", os.Getenv("AZURE_CLIENT_SECRET"))
	fmt.Println("AZURE_TENANT_ID:", os.Getenv("AZURE_TENANT_ID"))

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"subscription_id": os.Getenv("AZURE_SUBSCRIPTION_ID"),
			"client_id":       os.Getenv("AZURE_CLIENT_ID"),
			"client_secret":   os.Getenv("AZURE_CLIENT_SECRET"),
			"tenant_id":       os.Getenv("AZURE_TENANT_ID"),
			"postfix":         strings.ToLower(uniquePostfix),
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
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
	storageAccountExists := azure.StorageAccountExists(t, storageAccountName, resourceGroupName, subscriptionID)
	assert.True(t, storageAccountExists, "storage account does not exist")

	containerExists := azure.StorageBlobContainerExists(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscriptionID)
	assert.True(t, containerExists, "storage container does not exist")

	fileShareExists := azure.StorageFileShareExists(t, storageFileShareName, storageAccountName, resourceGroupName, subscriptionID)
	assert.True(t, fileShareExists, "File share does not exist")

	publicAccess := azure.GetStorageBlobContainerPublicAccess(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscriptionID)
	assert.False(t, publicAccess, "storage container has public access")

	accountKind := azure.GetStorageAccountKind(t, storageAccountName, resourceGroupName, subscriptionID)
	assert.Equal(t, storageAccountKind, accountKind, "storage account kind mismatch")

	skuTier := azure.GetStorageAccountSkuTier(t, storageAccountName, resourceGroupName, subscriptionID)
	assert.Equal(t, storageAccountTier, skuTier, "sku tier mismatch")

	actualDNSString := azure.GetStorageDNSString(t, storageAccountName, resourceGroupName, subscriptionID)
	storageSuffix, _ := azure.GetStorageURISuffixE()
	expectedDNS := fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix)
	assert.Equal(t, expectedDNS, actualDNSString, "Storage DNS string mismatch")
}
