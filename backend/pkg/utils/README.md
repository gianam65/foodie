# 🔧 Utils Package

Package `utils` chứa các helper functions dùng chung trong project.

## 📁 Cấu Trúc

```
pkg/utils/
├── http/          # HTTP response helpers
├── validation/    # Validation helpers
├── string/        # String manipulation utilities
├── time/          # Time formatting and utilities
├── id/            # ID generation and validation
└── pagination/    # Pagination helpers
```

## 📦 Packages

### `pkg/utils/http` - HTTP Response Helpers

Helper functions để gửi HTTP responses:

```go
import httputils "foodie/backend/pkg/utils/http"

// Success response (200 OK)
httputils.Success(w, data)

// Created response (201 Created)
httputils.Created(w, data)

// Error responses
httputils.BadRequest(w, "Invalid input", err)
httputils.NotFound(w, "Resource not found")
httputils.InternalServerError(w, "Server error", err)
httputils.Unauthorized(w, "Unauthorized")
httputils.Forbidden(w, "Forbidden")

// Custom status code
httputils.JSON(w, http.StatusAccepted, data)
httputils.Error(w, http.StatusConflict, "Conflict", err)
```

### `pkg/utils/validation` - Validation Helpers

Validation functions:

```go
import validation "foodie/backend/pkg/utils/validation"

// String validation
if validation.IsEmpty(email) { ... }
if validation.IsValidEmail(email) { ... }
if validation.IsValidUUID(id) { ... }

// Number validation
if err := validation.ValidateRange(age, 0, 120); err != nil { ... }
if err := validation.ValidatePositive(quantity); err != nil { ... }

// Slice validation
if err := validation.ValidateNonEmpty(items, "items"); err != nil { ... }
```

### `pkg/utils/string` - String Utilities

String manipulation functions:

```go
import stringutils "foodie/backend/pkg/utils/string"

// Truncate
short := stringutils.Truncate(longString, 50)

// Contains check
if stringutils.ContainsAny(text, "error", "fail", "exception") { ... }

// Case conversion
camel := stringutils.ToCamelCase("hello_world")  // "helloWorld"
snake := stringutils.ToSnakeCase("HelloWorld")   // "hello_world"

// Sanitize
clean := stringutils.Sanitize(userInput)

// Default value
name := stringutils.Default(userName, "Anonymous")
```

### `pkg/utils/time` - Time Utilities

Time formatting and manipulation:

```go
import timeutils "foodie/backend/pkg/utils/time"

// Formatting
formatted := timeutils.FormatRFC3339(time.Now())

// Parsing
t, err := timeutils.ParseRFC3339("2024-12-05T14:30:00Z")

// Timestamps
unix := timeutils.UnixTimestamp(time.Now())
unixMillis := timeutils.UnixTimestampMillis(time.Now())

// Time manipulation
future := timeutils.AddDays(time.Now(), 7)
future := timeutils.AddHours(time.Now(), 2)
future := timeutils.AddMinutes(time.Now(), 30)

// Duration
duration := timeutils.DurationBetween(startTime, endTime)
```

### `pkg/utils/id` - ID Generation

ID generation and validation:

```go
import idutils "foodie/backend/pkg/utils/id"

// UUID
uuid := idutils.GenerateUUID()
if idutils.IsValidUUID(uuid) { ... }

// Short ID (16 characters)
shortID, err := idutils.GenerateShortID()
shortID := idutils.MustGenerateShortID() // panics on error

// Sanitize
cleanID := idutils.SanitizeID(userInput)
```

### `pkg/utils/pagination` - Pagination Helpers

Pagination utilities:

```go
import pagination "foodie/backend/pkg/utils/pagination"

// Parse pagination params
page := pagination.ParsePage(pageStr, 1)
limit := pagination.ParseLimit(limitStr, 20, 1, 100)

// Calculate offset
offset := pagination.CalculateOffset(page, limit)

// Validate
if err := pagination.ValidatePagination(page, limit, 100); err != nil { ... }
```

## 💡 Ví Dụ Sử Dụng

### HTTP Controller

```go
import (
    httputils "foodie/backend/pkg/utils/http"
    validation "foodie/backend/pkg/utils/validation"
)

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
    // Validate
    if validation.IsEmpty(req.Email) {
        httputils.BadRequest(w, "Email is required", nil)
        return
    }

    // Process...

    httputils.Created(w, result)
}
```

### Use Case với Validation

```go
import (
    validation "foodie/backend/pkg/utils/validation"
    idutils "foodie/backend/pkg/utils/id"
)

func (uc *UseCase) Create(req CreateRequest) error {
    if !validation.IsValidEmail(req.Email) {
        return fmt.Errorf("invalid email")
    }

    id := idutils.GenerateUUID()
    // ...
}
```

## 🎯 Best Practices

1. **Import với alias**: Dùng alias để tránh conflict

   ```go
   import httputils "foodie/backend/pkg/utils/http"
   ```

2. **Organize theo domain**: Mỗi sub-package có một domain cụ thể

   - `http/` - HTTP related
   - `validation/` - Validation
   - `string/` - String operations
   - etc.

3. **Keep it simple**: Utils should be simple, pure functions

   - No dependencies on domain/infrastructure layers
   - Easy to test
   - Reusable across the project

4. **Document well**: Export functions should have clear documentation
