package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type StringSlice []string

func (s StringSlice) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		str := string(v)
		if str == "" {
			s = []string{}
		} else {
			s = strings.Split(str, ",")
		}
	case string:
		if v == "" {
			s = []string{}
		} else {
			s = strings.Split(v, ",")
		}
	case nil:
		s = []string{}
	default:
		return fmt.Errorf("unsupported scan, storing driver.Value type %T into type *[]string", value)
	}

	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "", nil
	}
	return strings.Join(s, ","), nil
}
