# Azure provider configuration
provider "azurerm" {
  features {}
  subscription_id = var.subscription
}

terraform {
  required_version = ">= 0.12.26"
  required_providers {
    azurerm = {
      version = "~> 2.20"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "resource_group" {
  name     = "terratest-storage-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A STORAGE ACCOUNT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_storage_account" "storage_account" {
  name                     = "storage${var.postfix}"
  resource_group_name      = azurerm_resource_group.resource_group.name
  location                 = azurerm_resource_group.resource_group.location
  account_kind             = var.storage_account_kind
  account_tier             = var.storage_account_tier
  account_replication_type = var.storage_replication_type
}

# Blob Container
resource "azurerm_storage_container" "container" {
  name                  = "blobcontainer1"
  storage_account_name  = azurerm_storage_account.storage_account.name
  container_access_type = var.container_access_type
}

# File Share
resource "azurerm_storage_share" "myfileshare" {
  name                 = "myfileshare"
  storage_account_name = azurerm_storage_account.storage_account.name
  quota                = 10 # Set the quota in GB
}