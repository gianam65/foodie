package dbtypes

// This package is kept for backward compatibility.
// The application now only uses SQL database (PostgreSQL/MySQL).
// MongoDB support has been removed.

// Type represents the backing database used by the application.
// Currently only SQL databases are supported.
type Type string

const (
	// TypeSQL represents SQL databases (PostgreSQL, MySQL)
	TypeSQL Type = "sql"
)
