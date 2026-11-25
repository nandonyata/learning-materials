terraform { // this "terraform" block defines the provider and Terraform versions you will use with this project
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
provider "aws" { // each provider block configures a specific provider, in this case, the AWS provider. The configuration includes the region, access key, secret key, and endpoint settings for LocalStack.
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

resource "aws_s3_bucket" "my_bucket" { // resource blocks define the infrastructure resources you want to create. In this case, it defines an S3 bucket resource with a specified bucket prefix and a force_destroy option that allows the bucket to be deleted even if it contains objects.
  bucket_prefix = "mybucketprefix-"
  force_destroy = true
}

output "bucket_name" {
  value = aws_s3_bucket.my_bucket.bucket
}