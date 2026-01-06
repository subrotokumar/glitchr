output "storage_trigger_queue_arn" {
  value = aws_sqs_queue.storage_trigger_queue.arn
}

output "storage_trigger_queue_url" {
  value = aws_sqs_queue.storage_trigger_queue.url
}