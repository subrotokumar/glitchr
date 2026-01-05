# Overview

This system is a **cloud-native adaptive bitrate (ABR) video streaming platform** that supports secure uploads, asynchronous multi-quality transcoding, and low-latency global playback.

The platform is designed for **production-scale workloads**, separating synchronous API traffic from long-running video processing while remaining cost-efficient and horizontally scalable.

## Technology Stack

* **Backend:** Go
* **Mobile Client:** Flutter
* **Container & Orchestration:** Docker, ECR, EKS
* **Authentication:** Amazon Cognito
* **Messaging:** Amazon SQS
* **Storage:** Amazon S3
* **CDN:** Amazon CloudFront
* **Databases:** PostgreSQL (RDS), Redis
* **Infrastructure:** Terraform / OpenTofu