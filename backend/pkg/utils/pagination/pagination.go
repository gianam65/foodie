package pagination

import (
	"fmt"
	"strconv"
)

// Meta contains pagination metadata for API responses.
type Meta struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"` // limit
	Offset      int  `json:"offset"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

// CalculateMeta calculates pagination metadata from total count, page, and limit.
func CalculateMeta(page, limit, total int) Meta {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit
	totalPages := (total + limit - 1) / limit // Ceiling division

	return Meta{
		CurrentPage: page,
		PerPage:     limit,
		Offset:      offset,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}

// ParsePage parses page number from query parameter, defaulting to 1.
func ParsePage(pageStr string) int {
	if pageStr == "" {
		return 1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// ParseOffset parses offset from query parameter.
func ParseOffset(offsetStr string) int {
	if offsetStr == "" {
		return 0
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}

// ParseLimit parses limit from query parameter, with min/max bounds.
func ParseLimit(limitStr string, defaultLimit, minLimit, maxLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return defaultLimit
	}

	if limit < minLimit {
		return minLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// CalculateOffset calculates the offset from page and limit.
func CalculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

// CalculatePageFromOffset calculates page number from offset and limit.
func CalculatePageFromOffset(offset, limit int) int {
	if limit <= 0 {
		return 1
	}
	page := (offset / limit) + 1
	if page < 1 {
		return 1
	}
	return page
}

// ValidatePagination validates page, offset, and limit values.
func ValidatePagination(page, offset, limit, maxLimit int) error {
	if page < 1 {
		return fmt.Errorf("page must be >= 1, got %d", page)
	}
	if offset < 0 {
		return fmt.Errorf("offset must be >= 0, got %d", offset)
	}
	if limit < 1 {
		return fmt.Errorf("limit must be >= 1, got %d", limit)
	}
	if limit > maxLimit {
		return fmt.Errorf("limit must be <= %d, got %d", maxLimit, limit)
	}
	return nil
}
