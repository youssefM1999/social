# Changelog

All notable changes to this project will be documented in this file.

## [1.2.0](https://github.com/youssefM1999/social/compare/v1.1.0...v1.2.0) (2026-01-02)


### Features

* update api version automatically ([0196178](https://github.com/youssefM1999/social/commit/019617894870949a5f7ec4357832eb584ff73df8))

## [1.1.0](https://github.com/youssefM1999/social/compare/v1.0.0...v1.1.0) (2026-01-02)


### Features

* add automation workflow ([026f8bc](https://github.com/youssefM1999/social/commit/026f8bcaaaf717a97f848dd280344a9f6df5d475))
* my release please script ([8614a9f](https://github.com/youssefM1999/social/commit/8614a9f272434030f52ba160187fae07e46fdcf5))


### Bug Fixes

* the deployment yaml file ([76e86b5](https://github.com/youssefM1999/social/commit/76e86b5fabaac7c949a58dd804f3e629e3100fc9))

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
