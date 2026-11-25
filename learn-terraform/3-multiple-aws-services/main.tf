provider "aws" {
  region     = "us-east-1"
  access_key = "test"
  secret_key = "test"

  skip_credentials_validation = true
  skip_requesting_account_id  = true
  s3_use_path_style           = true   # LocalStack needs path-style S3 URLs

  # Without an endpoints block, Terraform would send the requests to real AWS (and fail with the fake test credentials).
  endpoints {
    s3       = "http://localhost:4566"
    dynamodb = "http://localhost:4566"
  }
}

resource "aws_s3_bucket" "example" {
  bucket = "my-golang-app-bucket"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}

output "s3_bucket_name" {
  value = aws_s3_bucket.example.bucket
}

# -------------

resource "aws_dynamodb_table" "go-app-file" {
  name         = "go-app-file"
  billing_mode = "PAY_PER_REQUEST"   # no capacity planning needed
  hash_key     = "ID" # Go code in "./backend/main.go" reads and writes the key as ID

  attribute {
    name = "ID"
    type = "S"                       # S = string
  }
}

output "table_name" {
  value = aws_dynamodb_table.go-app-file.name
}