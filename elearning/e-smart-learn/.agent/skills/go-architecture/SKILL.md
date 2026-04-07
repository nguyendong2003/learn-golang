---
name: go-architecture
description: Overall directory structure, file naming conventions, and DTO usage in the e-Smart Learn project.
---

# Go Architecture & Project Convention Guidelines

This skill provides a high-level view of the e-Smart Learn API structure and conventions.

## 1. Directory Structure (`app/`)
- `apperror/`: Custom application-level error definitions.
- `cmd/`: Application entrypoints (e.g., `main.go`, CLI tooling setup).
- `config/`: Configurations struct mappings and env loaders.
- `consts/`: Constants used throughout the application (e.g., roles or statuses).
- `dto/`: Data Transfer Objects for API requests and responses. Separation from DB models.
- `handler/`: Gin web handlers, routing logic, and HTTP parsing.
- `job/` / `worker/`: Background jobs and task execution definitions.
- `model/`: GORM database structs.
- `pkg/`: External reusable packages (e.g., `storageProvider` wrappers).
- `repository/`: Database access mechanisms using GORM.
- `service/`: Core business logic connecting handlers to repositories.
- `util/`: Helper functions (validation logic, extracting from contexts, logging).

## 2. File Naming Conventions
- Go files follow `snake_case`. For instance, `user_service.go`, `auth_handler.go`.
- Corresponding test files if applicable would be `[file]_test.go`.
- Avoid stuttering inside subpackages (e.g., inside package `dto`, name files relevantly without full redundancy).

## 3. The DTO Pattern
The boundaries between the Model and the Client are strictly defined by DTOs in `app/dto/`.
- Handlers parse client JSON directly into Requests (e.g., `UpdateUserRequest`).
- Services perform logic referencing DB models.
- Services return Responses (`UserResponse`, `UserListResponse`). Models are rarely leaked to the UI layer.
- Ensure all struct fields heavily incorporate JSON tags (`json:"avatar"`).

## 4. Initialization & Overarching DI
The application uses clear Dependency Injection wire-ups generally initiated in `cmd/`.
Layers are pieced together hierarchically:
`DbRepository` -> `SpecificEntityRepository` -> `SpecificEntityService` -> `SpecificEntityHandler`.
No singletons or global variables are used for layer orchestration.
