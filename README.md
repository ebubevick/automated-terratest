
# Terratest for Azure Infrastructure

This repository demonstrates how to use **Terratest** to test the deployment of Azure resources using Terraform. Specifically, it deploys an Azure Resource Group, Storage Account, Blob Container, and File Share, then verifies that these resources exist and have the correct properties.

## Overview

- **Terraform**: Infrastructure-as-Code tool used to deploy resources to Azure.
- **Terratest**: A Go library for writing automated tests for Terraform code.
- **Azure**: Cloud provider where resources are deployed (e.g., Resource Group, Storage Account, Blob Container, and File Share).

This repository consists of the following components:
1. **Terraform Configuration**: Code that defines the Azure infrastructure.
2. **Terratest Go Tests**: Automated tests that verify the deployment and properties of Azure resources created by Terraform.


## Directory Structure

```plaintext
.
├── main.tf              # Terraform code to define resources (Resource Group, Storage Account, etc.)
├── outputs.tf           # Output values from Terraform (e.g., resource names)
├── terraform.tfstate     # State file for Terraform
├── terraform.tfstate.backup # Backup of the state file
├── terratest
│   ├── go.mod           # Go module definition
│   ├── go.sum           # Go module checksum file
│   └── main_test.go     # Terratest code to test the Terraform resources
└── variables.tf         # Variables for Terraform configuration
```

## Resources Deployed

The following resources are deployed using Terraform:
1. **Resource Group**: A container for Azure resources.
   - Name: `terratest-storage-rg-${var.postfix}`
2. **Storage Account**: A storage account to store blobs and files.
   - Name: `storage${var.postfix}`
   - Kind: Configured via the `var.storage_account_kind` variable.
   - Tier: Configured via the `var.storage_account_tier` variable.
3. **Blob Container**: A container within the Storage Account for storing blobs.
   - Name: `blobcontainer1`
4. **File Share**: A file share within the Storage Account for file storage.
   - Name: `myfileshare`
   - Quota: 10 GB

## Prerequisites

Before running the tests, ensure that you have the following tools installed:
- **Terraform**: Version 0.12.26 or later
- **Go**: Go 1.16 or later
- **Azure CLI**: For managing Azure resources from the command line

### Azure Authentication

Make sure you authenticate to Azure using one of the following methods:
- **Service Principal**: Set the following environment variables:
  ```bash
  export ARM_CLIENT_ID="your-client-id"
  export ARM_CLIENT_SECRET="your-client-secret"
  export ARM_SUBSCRIPTION_ID="your-subscription-id"
  export ARM_TENANT_ID="your-tenant-id"
  ```
- **Azure CLI**: Use `az login` to authenticate via Azure CLI.

## Setting Up the Environment

1. **Clone this repository**:
   ```bash
   git clone https://github.com/ebubevick/automated-terratest.git
   cd terraform-azure-terratest
   ```

2. **Initialize Terraform**:
   ```bash
   terraform init
   ```

3. **Configure the Variables**:
   Ensure that the `variables.tf` file is configured with your desired values. The most important variables are:
   - `subscription`: Your Azure subscription ID.
   - `location`: The Azure region for your resources.
   - `storage_account_kind`, `storage_account_tier`, `storage_replication_type`, and `container_access_type`.

4. **Apply Terraform**:
   Deploy the resources to Azure by running:
   ```bash
   terraform apply
   ```

5. **Run the Terratest**:
   After deploying the resources, you can run the Go tests to verify that the resources were created correctly:
   ```bash
   cd terratest
   go test -v
   ```

## Terratest Code Explanation

The `main_test.go` file contains the Terratest code that:
1. Initializes Terraform and applies the configuration.
2. Verifies the existence and properties of the following Azure resources:
   - Storage Account
   - Blob Container
   - File Share
3. Ensures the resources have the correct settings (e.g., no public access for the container, matching SKU tier, correct DNS for the storage account).

### Key Assertions:
- **Storage Account Exists**: Ensures the storage account exists in Azure.
- **Blob Container Exists**: Ensures the blob container exists within the storage account.
- **Storage Container Public Access**: Verifies that the blob container does not have public access.
- **Storage Account Kind and Tier**: Ensures the correct kind and tier are used for the storage account.
- **Storage DNS**: Verifies the DNS string for accessing the storage account.

## Cleanup

To remove the deployed resources and avoid unnecessary charges, run:
```bash
terraform destroy
```

## Conclusion

This repository demonstrates how to use Terratest to validate the creation and configuration of Azure resources deployed with Terraform. It tests the correctness of the resource properties and ensures the desired state is achieved in the Azure environment.