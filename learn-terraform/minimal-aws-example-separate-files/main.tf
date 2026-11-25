provider "aws" {
  region                      = var.region
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_requesting_account_id  = true

  endpoints {
    s3  = "http://172.24.34.66:4566"
    iam = "http://172.24.34.66:4566"
  }
}

resource "aws_s3_bucket" "test" {
  bucket_prefix = "${var.project_name}-"
  force_destroy = true
}
