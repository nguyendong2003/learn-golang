---
name: go-business-logic
description: Guidelines for service layer patterns, dependency injection, and error handling in the e-Smart Learn project.
---

# Go Business Logic Guidelines

This skill focuses on the service layer, where business logic is executed in the e-Smart Learn Go backend.

## 1. Service Layer Structure
Services are located in `app/service/` and named `[entity]_service.go`.
- Define an interface exposing the service methods (e.g., `UserService`).
- Create an unexported struct implementing the interface.
- Provide a `New[Entity]Service` constructor for dependency injection.

## 2. Dependency Injection
- Constructors receive instances of repositories, other services, or utility providers (e.g., `pkg.StorageProvider`).
- Service logic should not be coupled with HTTP logic (`gin.Context`); instead, pass standard `context.Context` from the handler.
```go
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	db             repository.DbRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	db repository.DbRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
		db:             db,
	}
}
```

## 3. Error Handling
- Use the custom `apperror` package to return classified errors (`apperror.NewNotFoundError`, `apperror.NewBadRequestError`, `apperror.NewInternalServerError`, etc.).
- Never panic. Return descriptive error wrappers instead of generic database errors directly to the client.
```go
user, err := s.userRepository.FindByID(ctx, id, nil)
if err != nil {
	return nil, apperror.NewInternalServerError("Failed to get user detail")
}
if user == nil {
	return nil, apperror.NewNotFoundError("User not found")
}
```

## 4. Transaction Management
- When involving multiple database changes (e.g., creating a user and their profile), use transactions via the `DbRepository` which provides `.Transaction()` or `.Begin()`, `.Commit()`, and `.Rollback()`.
```go
err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
	// Execute queries with transaction repositories
	return nil
})
```

## 5. DTO Mapping
Services should accept primitives or DTO structures (`dto.UpdateUserRequest`) and return DTO response structures (`dto.UserResponse`), NOT raw GORM models unless strictly for internal use. Map domain models to DTOs before returning.
