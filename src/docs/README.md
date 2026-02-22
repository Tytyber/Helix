# API Service Documentation

This document provides an overview of the API service structure, its components, and how they interact.

## Overview

The API is a Go-based backend service built using the [Gin](https://gin-gonic.com/) framework for routing and PostgreSQL via [pgx](https://github.com/jackc/pgx) for database access. It implements standard RESTful endpoints for user authentication, social features (friends, groups, hubs), real-time communication (chat, calls), and content feeds.

The application follows a clean architecture with separation between:
- **Handlers** (`hundler.go`): HTTP request/response logic
- **Services** (`service.go`): Business logic and database interactions
- **Models/Types**: Data structures used across layers

Authentication is handled via JWT tokens, with refresh tokens stored securely in the database.

## Project Structure

```
api/
├── go.mod
├── go.sum
├── main.go
└── internal/
    ├── auth/           # Authentication & user management
    ├── calls/          # Voice/video calls
    ├── chat/           # Direct messaging
    ├── config/         # Configuration loading
    ├── database/       # Database connection
    ├── feed/           # Personalized content feed
    ├── friend/         # Friend system
    ├── group/          # Group chats
    ├── hub/            # Community hubs with posts
    └── middleware/     # Shared middleware (e.g. auth)
```

## Core Components

### Main Application (main.go)

The `main.go` file initializes:
- Configuration via environment variables
- Database connection pool
- All service instances
- HTTP routes grouped by functionality

Routes are organized under `/api/v1` with:
- Public routes under `/auth`
- Protected routes (JWT required) under `/protected`

### Configuration (internal/config)

Uses `godotenv` to load environment variables with fallbacks:
- `PORT` (default: ":8080")
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- `JWT_SECRET` (default: "supersecret")

### Database (internal/database)

Manages a connection pool using `pgxpool`. The pool is passed to all services that require database access.

### Authentication (internal/auth)

#### Endpoints
- `POST /api/v1/auth/register` — Register new user
- `POST /api/v1/auth/login` — Login and get tokens
- `POST /api/v1/auth/refresh` — Refresh JWT token
- `POST /api/v1/auth/logout` — Invalidate refresh token
- `GET /api/v1/protected/me` — Get current user profile
- `PUT /api/v1/protected/profile` — Update profile
- `POST /api/v1/protected/reputation` — Add reputation to another user

#### Implementation
- Passwords hashed with `bcrypt`
- JWT access tokens (15 min expiry)
- Refresh tokens stored as SHA256 hashes in DB
- Middleware validates `Bearer` tokens

### Hubs (internal/hub)

Community spaces where users can post content.

#### Endpoints
- `GET /api/v1/protected/hubs` — List all hubs
- `GET /api/v1/protected/hubs/:id` — Get specific hub
- `POST /api/v1/protected/hubs/:id/join` — Join hub
- `POST /api/v1/protected/hubs/:id/leave` — Leave hub
- `GET /api/v1/protected/hubs/:id/posts` — Get hub posts
- `POST /api/v1/protected/hubs/:id/posts` — Create post

### Feed (internal/feed)

Personalized content feed based on user's hubs.

#### Endpoints
- `GET /api/v1/protected/feed` — Get user's feed
- `GET /api/v1/protected/feed/top-hubs` — Get top hubs by membership

### Friends (internal/friend)

Friend requests and social graph.

#### Endpoints
- `POST /api/v1/protected/friends/request` — Send friend request
- `POST /api/v1/protected/friends/accept` — Accept request
- `GET /api/v1/protected/friends` — List friends

### Groups (internal/group)

Private group chats with multiple participants.

#### Endpoints
- `POST /api/v1/protected/groups` — Create group
- `POST /api/v1/protected/groups/:id/join` — Join group
- `POST /api/v1/protected/groups/:id/leave` — Leave group
- `GET /api/v1/protected/groups/:id/messages` — Get messages
- `POST /api/v1/protected/groups/:id/messages` — Send message

### Chat (internal/chat)

Direct 1-on-1 messaging.

#### Endpoints
- `POST /api/v1/protected/chats` — Create chat with user
- `GET /api/v1/protected/chats/:id/messages` — Get chat history
- `POST /api/v1/protected/chats/:id/messages` — Send message

### Calls (internal/calls)

Voice/video call management.
protected.POST("/friends/request", friendsHandler.SendRequest)
protected.POST("/friends/accept", friendsHandler.AcceptRequest)
protected.GET("/friends", friendsHandler.ListFriends)

protected.POST("/calls", callsHandler.CreateCall)
protected.PUT("/calls/:id/status", callsHandler.UpdateCallStatus)
#### Endpoints
- `POST /api/v1/protected/calls` — Initiate call
- `PUT /api/v1/protected/calls/:id/status` — Update call status (pending, accepted, rejected, ended)

## Data Models

### User
- id, username, email, password_hash, avatar_url, bio, status, reputation, created_at

### Hub
- id, name, description, avatar, created_at

### HubPost
- id, hub_id, author_id, content, likes, created_at

### FeedPost (read-only)
- id, hub_id, hub_name, author_id, content, likes, created_at

### TopHub (read-only)
- id, name, description, avatar, members

### Group
- id, name, owner_id, created_at

### GroupMessage
- id, group_id, sender_id, content, created_at

### Chat
- id, user1_id, user2_id, created_at

### Message
- id, chat_id, sender_id, content, created_at

### Call
- id, caller_id, callee_id, status, created_at

## Error Handling

Most endpoints return JSON errors:
- `400 Bad Request` — Invalid input
- `401 Unauthorized` — Missing/invalid token
- `404 Not Found` — Resource not found
- `500 Internal Server Error` — Server-side error

## Dependencies

See `go.mod` for full list. Key dependencies:
- `github.com/gin-gonic/gin` — Web framework
- `github.com/jackc/pgx/v5` — PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` — JWT implementation
- `golang.org/x/crypto/bcrypt` — Password hashing
- `github.com/joho/godotenv` — Environment variables

## How to Run

1. Copy `.env.example` to `.env` and configure
2. Ensure PostgreSQL is running
3. Run `go run main.go`

## Future Improvements

- Add input validation
- Implement rate limiting
- Add logging middleware
- Support pagination in feed and messages
- Add unit and integration tests
- Implement WebSocket support for real-time updates
