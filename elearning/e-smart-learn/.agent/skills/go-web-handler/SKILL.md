---
name: go-web-handler
description: Guidelines for creating Gin web handlers, routing, and JSON responses in the e-Smart Learn project.
---

# Go Web Handler Guidelines

This skill focuses on how web handlers are structured and implemented in the e-Smart Learn Go backend project.

## 1. Handler Definition & Interface
Handlers are defined in `app/handler/` using the format `[entity]_handler.go`.
- Define an interface for the handler containing methods returning `gin.HandlerFunc`.
- Create an unexported struct implementing the interface, holding references to required services.
- Provide a `New[Entity]Handler` constructor function for dependency injection.

```go
type UserHandler interface {
	GetByID() gin.HandlerFunc
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}
```

## 2. Handling `gin.Context`
- Handlers return `gin.HandlerFunc` closures: `return func(c *gin.Context) { ... }`.
- Path parameters: `c.Param("id")`.
- Query parameters: Bind using `c.ShouldBindQuery(&request)`.
- Body payloads: Use custom utility `util.BindAndValidateJSON(c, &request)` for JSON body binding and validation.
- User ID (from auth middleware context): Use `util.GetRequestUserID(c)`.

## 3. Returning JSON Responses
Use the provided `dto.ApiResponse` helpers for standardizing responses:
- Success responses encapsulate data, metadata (like pagination), and request info:
```go
res := dto.NewApiResponse(c)
res.Request = dto.GetRequestClient(c)
res.Data = data
// For pagination:
res.Metadata = dto.NewPagination(limit, offset, int(total), "created_at", "desc")
c.JSON(http.StatusOK, res)
```

## 4. Error Handling in Handlers
- Pass errors down to Gin's error middleware using `_ = c.Error(err)`. Do NOT call `c.JSON` manually for errors. Wait and `return`.
- Example for invalid UUID:
```go
if err != nil {
	_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
	return
}
```
- Example for passing service layer errors:
```go
data, err := h.userService.GetByID(c.Request.Context(), userID)
if err != nil {
	_ = c.Error(err)
	return
}
```

## 5. Swagger Annotations
Add godoc annotations directly above each handler constructor or router function for Swagger generation:
```go
// @Summary Get user by ID
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// ...
// @Router /api/v1/admin/users/{id} [get]
```
