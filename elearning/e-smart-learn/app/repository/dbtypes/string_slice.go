package dbtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice []string

func (a StringSlice) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *StringSlice) Scan(src interface{}) error {
	if src == nil {
		*a = []string{}
		return nil
	}
	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, a)
	case string:
		return json.Unmarshal([]byte(s), a)
	default:
		return fmt.Errorf("cannot scan %T", src)
	}
}
