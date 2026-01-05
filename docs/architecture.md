
# System Architecture

## High-Level Architecture

```mermaid
flowchart LR
    Client[Flutter App]
    API[Go API Server]
    S3[(Amazon S3)]
    CDN[CloudFront]
    SQS[SQS Queue]
    Consumer[Go Consumer]
    RDS[(Postgres RDS)]
    Redis[(Redis Cache)]

    Client -->|Auth / API| API
    API -->|Metadata| RDS
    API -->|Cache| Redis
    API -->|Pre-signed URL| Client
    Client -->|Upload| S3
    S3 -->|Event| SQS
    SQS --> Consumer
    Consumer -->|Transcode| S3
    S3 --> CDN
    CDN --> Client
```

## Architectural Principles

* Stateless API services
* Event-driven async processing
* Independent scaling of API and consumers
* CDN-first content delivery