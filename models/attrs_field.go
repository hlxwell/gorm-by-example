package models

import "fmt"

type AttrsField struct {
	string
}

func (c AttrsField) MarshalJSON() ([]byte, error) {
	return []byte(c.Value()), nil
}

func (c AttrsField) Value() string {
	return fmt.Sprintf("\"%s\"", c.string)
}

func (c *AttrsField) Scan(v interface{}) error {
	value, ok := v.(string)
	if ok {
		*c = AttrsField{value}
		return nil
	}

	return fmt.Errorf("can not convert %v to string", v)
}
