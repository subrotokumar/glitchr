# Backend Services (Go)

## API Server

### Responsibilities

* Authentication and authorization
* Upload orchestration
* Video metadata APIs
* Playback manifest exposure

### Key Endpoints

| Method | Endpoint     | Description           |
| ------ | ------------ | --------------------- |
| POST   | /videos      | Create upload session |
| GET    | /videos/{id} | Video metadata        |
| GET    | /health      | Liveness / readiness  |

### Design Notes

* Stateless handlers
* Token validation middleware (Cognito)
* Redis-backed caching layer

## Consumer Service

### Responsibilities

* Poll SQS for jobs
* Execute FFmpeg-based transcoding
* Upload renditions and manifests

### Scaling Strategy

* Scale on SQS queue depth
* Independent deployment from API