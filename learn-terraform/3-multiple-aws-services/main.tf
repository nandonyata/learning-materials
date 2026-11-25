provider "aws" {
  region     = "us-east-1"
  access_key = "test"
  secret_key = "test"

  skip_credentials_validation = true
  skip_requesting_account_id  = true
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
  hash_key     = "uuid"

  attribute {
    name = "uuid"
    type = "S"                       # S = string
  }
}

output "table_name" {
  value = aws_dynamodb_table.go-app-file.name
}