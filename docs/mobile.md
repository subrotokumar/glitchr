# Flutter

## Responsibilities

* User authentication
* Video upload initiation
* Adaptive playback

## Upload Flow

```mermaid
sequenceDiagram
    participant U as User
    participant A as Flutter App
    participant API as Go API
    participant S3 as S3

    U->>A: Select Video
    A->>API: Request Upload URL
    API-->>A: Pre-signed URL
    A->>S3: Upload Video
```

## Playback

* Uses HLS master playlist
* Player automatically switches bitrate based on network conditions