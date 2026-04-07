---
name: go-db-repository
description: Guidelines for GORM/SQL models and repository patterns in the e-Smart Learn project.
---

# Go DB Repository Guidelines

This skill details how database models and data access using the repository pattern are implemented in the e-Smart Learn project.

## 1. GORM Models
Models are placed in `app/model/` and are named simply as `[entity].go`.
- Embed `BaseModel` (which includes `ID` as UUID, `CreatedAt`, `UpdatedAt`, `DeletedAt` for soft deletes).
- Define associations using GORM tags, maintaining clear foreign keys.
```go
type User struct {
	BaseModel
	Name     string    `gorm:"type:varchar(100);not null"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	RoleID   uuid.UUID `gorm:"type:uuid"`
	Role     *Role     `gorm:"foreignKey:RoleID"`
}
```

## 2. Repository Pattern Setup
Repositories abstract the database access. They sit in `app/repository/` and follow `[entity]_repository.go`.
- A generic `BaseRepository[T]` provides standard CRUD (`Create`, `FindByID`, `Update`, `Delete`, `Count`, etc.).
- Entity repositories embed the `Repository[T]` interface and `baseRepository` to inherit generic methods, then add entity-specific queries.

```go
type UserRepository interface {
	Repository[model.User] // Inherit generic methods
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	*repository[model.User]
}

func NewUserRepository(db DbRepository) UserRepository {
	return &userRepository{
		repository: NewBaseRepository[model.User](db), // initialize generic behaviors
	}
}
```

## 3. Querying & Preloading
- Define custom queries by building off a base query or utilizing Raw SQL via `.Raw()`.
- Return `*model.Entity` for found, `nil, nil` for not found (`gorm.ErrRecordNotFound`), and `nil, err` on failure.
- Preloads can be dynamically passed via the `repository.Preload` type system.

## 4. Complex Reads (CTEs and Stats)
For complex aggregates (like pagination with aggregated stats), the project relies heavily on raw Postgres queries (`WITH` CTEs) and struct mapping onto custom model rows (e.g., `model.UserDirectoryRow`).
```go
func (r *userRepository) GetList(...) ([]*model.UserDirectoryRow, int64, error) {
    // Write Raw SQL with CTEs counting lessons and mapping
    // db.Raw(query).Scan(&rows)
}
```
