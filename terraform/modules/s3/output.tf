output "raw_bucket_id" {
  value = aws_s3_bucket.raw_bucket.id
}

output "processed_bucket_id" {
  value = aws_s3_bucket.main_bucket.id
}