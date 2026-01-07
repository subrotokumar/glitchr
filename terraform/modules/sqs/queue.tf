resource "aws_sqs_queue" "sqs_queue" {
  name =  var.sqs_queue
}

resource "aws_sqs_queue_policy" "sqs_queue_policy" {
  queue_url = aws_sqs_queue.sqs_queue.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
        Action    = "sqs:SendMessage"
        Resource  = aws_sqs_queue.sqs_queue.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = var.bucket_arn
          }
        }
      }
    ]
  })
}