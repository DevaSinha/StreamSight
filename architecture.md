# Project Architecture: Go + Next.js Real-Time Dashboard

## High-Level Overview

+-----------------+ +-----------------+ +---------------+
| Next.js App |<---->| Go API Service |<---->| PostgreSQL |
| (Frontend+SSR) | +-----------------+ | Database |
| | ^ ^ +---------------+
| | | |
| [WebSocket/ | +--------------------+
| REST Client] |<---->| Go Worker/Stream |
| | | (Streaming/CV) |
+-----------------+ +--------------------+

text

## Major Blocks
- **Next.js App**: UI, routing, SSR, authentication, dashboard, video feeds, overlays, alerts, analytics
- **Go API Service**: Core REST/WebSocket API, auth, camera/user/event management
- **Go Worker**: Handles RTSP/mock video ingest, runs detection/analytics, delivers events
- **PostgreSQL**: Structured data for users, cameras, alerts, events
- **[Optional] Kafka/Redis**: For event/message queue decoupling if future scaling

## Component Breakdown

### Next.js Frontend
- JWT Auth, camera CRUD, live dashboard, WebRTC/WebSocket, overlays, analytics

### Go API Service
- REST endpoints (auth, cameras, alerts, analytics), WebSocket (alerts), JWT sessions, DB ORM

### Go Worker
- Handles all stream ingest and detection (face, emotion, mask, etc), posts events to API

### Database
- Tables: Users, Cameras, Alerts, Events, Overlays, (Roles if RBAC)

## Interactions Flow
1. Login via Next.js → `/api/auth/login` (to Go API)
2. Fetch dashboard/camera list via `/api/cameras`
3. Add camera via `/api/cameras`
4. Worker receives camera config, connects RTSP/mock
5. Detection event: POST `/api/detection-event` or push to queue
6. Go API: stores event, emits to `/ws/alerts`
7. Next.js: receives overlay, updates UI
8. Next.js: analytics GET `/api/analytics`

## Dev & Prod Tooling
- Docker Compose up: Next, Go API, Worker, Postgres
- GitHub Actions: build/test/release CI
- Demo: mock RTSP stream in Docker

## Extendibility
- Plugins for new models, RBAC, scalable queue, edge agent, K8s