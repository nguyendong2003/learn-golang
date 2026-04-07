---
name: go-aggregation-api
description: Pattern for building paginated, filterable aggregation/statistics APIs in the e-Smart Learn Go backend. Use when an endpoint must JOIN multiple tables, GROUP BY an entity, and return computed aggregates (sums, counts) with optional filters and pagination.
---

# Go Aggregation API Pattern

## When to Use

Apply this skill when the endpoint:
- Aggregates data **across multiple tables** (e.g. JOIN + GROUP BY)
- Returns a **list of entities** each with computed numeric metrics
- Supports **optional filters** (date range, status, role, etc.)
- Requires **pagination** (limit / offset)

Examples: revenue per instructor, enrollments per course, orders per user, active students per month.

---

## Architecture Layer Flow

```
HTTP Request
    ↓
Handler        — bind & validate query params → call service
    ↓
Service        — parse/validate filter values → call repo with primitives → map rows → DTOs
    ↓
Repository     — build dynamic SQL → scan into Row structs (cents/raw units)
    ↓
Database       — raw SQL CTE with optional WHERE clauses + LIMIT/OFFSET
```

---

## DTO Rules (`app/dto/`)

**Filter Request struct** — bound from HTTP query params:
```go
type XxxFilterRequest struct {
    StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
    EndDate   string `form:"end_date"   binding:"omitempty,datetime=2006-01-02"`
    Limit     int    `form:"limit"      binding:"omitempty,min=1,max=100"`
    Offset    int    `form:"offset"     binding:"omitempty,min=0"`
    SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

func (r *XxxFilterRequest) Process() {
    if r.Limit <= 0   { r.Limit = 10 }
    if r.Offset < 0   { r.Offset = 0 }
    if r.SortOrder == "" { r.SortOrder = "desc" }
}
```

**Item Response struct** — one record per aggregated entity:
```go
type XxxItemResponse struct {
    EntityID    string  `json:"entity_id"`
    EntityName  string  `json:"entity_name"`
    MetricA     float64 `json:"metric_a"`  // always dollars/display units, never raw cents
    MetricB     int64   `json:"metric_b"`
}
```

**Rules:**
- Request structs live in `dto/` — never in `repository/` or `service/`
- Response structs use display units (dollars, not cents; percentages as float64)
- No `ListResponse` wrapper struct needed — use `dto.NewPagination` in handler

---

## Repository Rules (`app/repository/`)

**Row struct** — raw DB scan target, raw units (cents, counts):
```go
// XxxRow is the raw DB scan target. Monetary fields are in cents.
type XxxRow struct {
    EntityID  string `gorm:"column:entity_id"`
    MetricA   int64  `gorm:"column:metric_a"`  // cents
    MetricB   int64  `gorm:"column:metric_b"`
}
```

**Interface method signature** — primitives only, NO structs:
```go
GetXxxAggregation(
    ctx       context.Context,
    startDate *time.Time,   // nil = no filter
    endDate   *time.Time,   // nil = no filter
    limit     int,
    offset    int,
    sortOrder string,
) ([]*XxxRow, int64, error)
```

**Rules:**
- Repository accepts **only primitive types** (`*time.Time`, `int`, `string`, `uuid.UUID`)
- **Never** accept a DTO or params struct in the repository layer
- Row struct stays in the `repository` package, not `model` (it's a query result, not a persisted entity)
- Return `nil, 0, err` on DB error — never wrap with `apperror` here (service's job)

---

## SQL Pattern

Use a **two-query approach**: count first, then list. Share the same CTE base.

```go
func (r *xxxRepository) GetXxxAggregation(
    ctx context.Context,
    startDate, endDate *time.Time,
    limit, offset int, sortOrder string,
) ([]*XxxRow, int64, error) {

    // 1. Build dynamic date filter
    var conditions []string
    var args []any
    if startDate != nil {
        conditions = append(conditions, "t.created_at >= ?")
        args = append(args, *startDate)
    }
    if endDate != nil {
        conditions = append(conditions, "t.created_at <= ?")
        args = append(args, *endDate)
    }
    dateFilter := ""
    if len(conditions) > 0 {
        dateFilter = "AND " + strings.Join(conditions, " AND ")
    }

    // 2. Shared CTE
    baseCTE := `
        WITH aggregated AS (
            SELECT entity_id, SUM(amount) AS metric_a
            FROM transactions t
            WHERE t.deleted_at IS NULL ` + dateFilter + `
            GROUP BY entity_id
        )
    `

    // 3. COUNT query (reuse same date args)
    var total int64
    countArgs := make([]any, len(args))
    copy(countArgs, args)
    countQuery := baseCTE + `SELECT COUNT(DISTINCT e.id) FROM entities e ...`
    if err := r.db.GetDB().WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
        return nil, 0, err
    }

    // 4. Sanitize sortOrder to prevent SQL injection
    if sortOrder != "asc" {
        sortOrder = "desc"
    }
    orderClause := fmt.Sprintf("metric_a %s", sortOrder)

    // 5. LIST query
    listArgs := append(args, limit, offset)
    listQuery := baseCTE + `
        SELECT e.id::text AS entity_id, COALESCE(a.metric_a, 0) AS metric_a
        FROM entities e
        LEFT JOIN aggregated a ON a.entity_id = e.id
        WHERE e.deleted_at IS NULL
        ORDER BY ` + orderClause + `
        LIMIT ? OFFSET ?
    `
    rows := make([]*XxxRow, 0)
    if err := r.db.GetDB().WithContext(ctx).Raw(listQuery, listArgs...).Scan(&rows).Error; err != nil {
        return nil, 0, err
    }
    return rows, total, nil
}
```

**SQL Rules:**
- Always use `COALESCE(col, 0)` for LEFT JOIN aggregates to avoid nil panics
- Always use `NULLIF(denominator, 0)` for divisions to prevent divide-by-zero
- Parameterize all user inputs via `?` — never interpolate directly
- `sortOrder` must be validated to `"asc"` or `"desc"` before interpolation into ORDER BY
- Soft-delete filter: always add `AND table.deleted_at IS NULL`

---

## Service Rules (`app/service/`)

```go
func (s *xxxService) GetXxxAggregation(
    ctx context.Context,
    req dto.XxxFilterRequest,
) ([]*dto.XxxItemResponse, int64, error) {

    // 1. Parse string dates → *time.Time (service responsibility, not repo)
    var startDate, endDate *time.Time
    if req.StartDate != "" {
        t, err := time.Parse("2006-01-02", req.StartDate)
        if err != nil {
            return nil, 0, apperror.NewBadRequestError("Invalid start_date, expected YYYY-MM-DD")
        }
        startDate = &t
    }
    if req.EndDate != "" {
        t, err := time.Parse("2006-01-02", req.EndDate)
        if err != nil {
            return nil, 0, apperror.NewBadRequestError("Invalid end_date, expected YYYY-MM-DD")
        }
        end := t.Add(24*time.Hour - time.Second) // inclusive end of day
        endDate = &end
    }

    // 2. Call repo with primitives
    rows, total, err := s.xxxRepository.GetXxxAggregation(ctx, startDate, endDate, req.Limit, req.Offset, req.SortOrder)
    if err != nil {
        return nil, 0, apperror.NewInternalServerError("Failed to retrieve aggregation data")
    }

    // 3. Map rows → response DTOs (convert units here, not in repo or handler)
    items := make([]*dto.XxxItemResponse, len(rows))
    for i, row := range rows {
        items[i] = &dto.XxxItemResponse{
            EntityID: row.EntityID,
            MetricA:  float64(row.MetricA) / 100.0, // cents → dollars
            MetricB:  row.MetricB,
        }
    }
    return items, total, nil
}
```

**Service Rules:**
- Service accepts DTO, passes **primitives** to repository
- Date parsing and validation happens in **service**, not handler or repo
- Unit conversion (cents → dollars, raw → percentage) happens in **service**, not repo
- Wrap all repo errors with `apperror.NewInternalServerError` — never leak raw DB errors

---

## Handler Rules (`app/handler/`)

```go
func (h *xxxHandler) GetXxxAggregation() gin.HandlerFunc {
    return func(c *gin.Context) {
        var req dto.XxxFilterRequest
        if err := c.ShouldBindQuery(&req); err != nil {
            _ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
            return
        }
        req.Process() // apply defaults

        data, total, err := h.xxxService.GetXxxAggregation(c.Request.Context(), req)
        if err != nil {
            _ = c.Error(err)
            return
        }

        res := dto.NewApiResponse(c)
        res.Request  = dto.GetRequestClient(c)
        res.Data     = data
        res.Metadata = dto.NewPagination(req.Limit, req.Offset, int(total), "metric_a", req.SortOrder)

        c.JSON(http.StatusOK, res)
    }
}
```

**Handler Rules:**
- Bind with `c.ShouldBindQuery` for GET; `util.BindAndValidateJSON` for POST/PUT
- Always call `req.Process()` after binding to apply defaults
- Errors go through `_ = c.Error(err)` then `return` — never call `c.JSON` for errors
- Handler does **no** business logic, **no** date parsing, **no** unit conversion

---

## Pagination Rules

- Always sort by SQL column name interpolated **after sanitization** (`asc`/`desc` check)
- `dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)` is the standard call
- Count query must **not** include `LIMIT`/`OFFSET`
- Use separate `countArgs` copy so list args can safely append `limit, offset`

---

## Anti-patterns ❌

| Anti-pattern | Why wrong |
|---|---|
| Passing a `QueryParams` struct to the repository | Couples repo to calling convention; breaks SRP |
| Parsing dates in the handler | Handler should only bind HTTP — parsing is service logic |
| Converting cents→dollars in the repository | Repo returns raw storage units; conversion is service logic |
| Interpolating `sortOrder` directly into SQL without validation | SQL injection risk |
| Using `apperror` inside the repository | Repo returns raw errors; service decides error classification |
| Using `LIMIT`/`OFFSET` in the count query | Produces wrong total count |
