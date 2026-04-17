package handlers

import (
	"database/sql"
	"html"
	"strings"
)

// helper function to wrap string to NullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// sanititze user input 
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = html.EscapeString(s)
	return s
}