variable "subscription_id" {
  description = "Azure subscription ID"
  type        = string
  default     = "4bf3e463-ed9f-4148-8906-3eed094e0794"
}

variable "client_id" {
  description = "Azure client ID"
  type        = string
}

variable "client_secret" {
  description = "Azure client secret"
  type        = string
  sensitive   = true
}

variable "tenant_id" {
  description = "Azure tenant ID"
  type        = string
}

variable "location" {
  description = "The location to set for the storage account."
  type        = string
  default     = "East US"
}

variable "storage_account_kind" {
  description = "The kind of storage account to set"
  type        = string
  default     = "StorageV2"
}

variable "storage_account_tier" {
  description = "The tier of storage account to set"
  type        = string
  default     = "Standard"
}

variable "storage_replication_type" {
  description = "The replication type of storage account to set"
  type        = string
  default     = "GRS"
}

variable "container_access_type" {
  description = "The access level for the storage container (e.g., private, blob, container)"
  type        = string
  default     = "private"
}

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions"
  type        = string
  default     = "resource"
}
