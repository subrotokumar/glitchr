resource "aws_s3_bucket" "raw_bucket" {
  bucket = var.raw_bucket
}

resource "aws_s3_bucket" "main_bucket" {
  bucket = var.main_bucket
}