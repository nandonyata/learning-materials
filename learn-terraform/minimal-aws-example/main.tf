terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}


# provider "aws" {
#     region = "us-east-1"
#     profile = "terraform_user"
# }
# 
# 
# For LocalStack, you don't need a profile — you need to configure the provider with endpoint and dummy credentials instead - https://docs.localstack.cloud/aws/tutorials/s3-static-website-terraform/
# Solution: Update your main.tf
# Replace the profile with explicit credentials and endpoint:
provider "aws" {
  region                      = "us-east-1"
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_requesting_account_id  = true

  endpoints {
    s3  = "http://172.24.34.66:4566"
    iam = "http://172.24.34.66:4566"
  }
}

resource "aws_s3_bucket" "my_bucket" {
  bucket_prefix = "mybucketprefix-"
  force_destroy = true
}

output "bucket_name" {
  value = aws_s3_bucket.my_bucket.bucket
}