# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2025-12-23

### Features

- **Authentication & Authorization**
  - Implement JWT authentication
  - Add authorization and role-based access control
  - Implement user validation with frontend and backend activation scheme
  - Add Basic Auth for protected endpoints

- **Caching**
  - Add Redis caching layer
  - Implement user caching for improved performance

- **API Endpoints**
  - Create posts API endpoint
  - Implement comments for posts endpoint
  - Create user feed endpoint
  - Create user storage with optimistic concurrency control
  - Implement follow/unfollow handling

- **Email**
  - Set up mailer using SendGrid
  - Implement user invitation emails

- **Infrastructure**
  - Implement structured logging with Zap
  - Add rate limiting middleware
  - Add server metrics endpoint
  - Implement graceful shutdown support
  - Add CI/CD automation workflow

### Architecture

- Established base API architecture
- Database migrations for users, posts, comments, roles, and followers
