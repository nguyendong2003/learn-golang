package dbtypes

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDSlice []uuid.UUID

func (a UUIDSlice) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(a))
	for i, u := range a {
		parts[i] = u.String()
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *UUIDSlice) Scan(src any) error {
	if src == nil {
		*a = UUIDSlice{}
		return nil
	}
	var s string
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("unsupported scan source type %T", src)
	}
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = UUIDSlice{}
		return nil
	}
	parts := strings.Split(s, ",")
	out := make(UUIDSlice, len(parts))
	for i, p := range parts {
		u, err := uuid.Parse(strings.TrimSpace(p))
		if err != nil {
			return err
		}
		out[i] = u
	}
	*a = out
	return nil
}

// Utility function to convert UUIDSlice to []string for easier JSON serialization
func ToStringSlice(a UUIDSlice) []string {
	ids := make([]string, len(a))
	for i, id := range a {
		ids[i] = id.String()
	}
	return ids
}
