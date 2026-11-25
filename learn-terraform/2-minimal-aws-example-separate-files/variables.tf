variable "region" { // use this later with "var.region" in the provider block
  type        = string
  description = "The AWS region to create resources in"
  default     = "us-east-1"
}

variable "project_name" {
  type        = string
  description = "Prefix for the S3 bucket and name"
  default     = "test-demo"
}