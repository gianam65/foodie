package time

import (
	"time"
)

// FormatRFC3339 formats a time.Time to RFC3339 string.
func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseRFC3339 parses an RFC3339 formatted string to time.Time.
func ParseRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// Now returns the current time.
func Now() time.Time {
	return time.Now()
}

// UnixTimestamp returns the Unix timestamp (seconds since epoch).
func UnixTimestamp(t time.Time) int64 {
	return t.Unix()
}

// UnixTimestampMillis returns the Unix timestamp in milliseconds.
func UnixTimestampMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// FromUnixTimestamp creates a time.Time from a Unix timestamp.
func FromUnixTimestamp(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}

// FromUnixTimestampMillis creates a time.Time from a Unix timestamp in milliseconds.
func FromUnixTimestampMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

// IsAfter checks if time t is after time u.
func IsAfter(t, u time.Time) bool {
	return t.After(u)
}

// IsBefore checks if time t is before time u.
func IsBefore(t, u time.Time) bool {
	return t.Before(u)
}

// AddDays adds n days to a time.
func AddDays(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n)
}

// AddHours adds n hours to a time.
func AddHours(t time.Time, n int) time.Time {
	return t.Add(time.Duration(n) * time.Hour)
}

// AddMinutes adds n minutes to a time.
func AddMinutes(t time.Time, n int) time.Time {
	return t.Add(time.Duration(n) * time.Minute)
}

// DurationBetween returns the duration between two times.
func DurationBetween(start, end time.Time) time.Duration {
	return end.Sub(start)
}
