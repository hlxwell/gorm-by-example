package models

import "fmt"

type UserAttributes struct {
	Name string   `json:"name"`
	Age  uint     `json:"age"`
	Tags []string `json:"tags"`
}

func (a UserAttributes) Validate() error {
	if a.Age <= 0 {
		return fmt.Errorf("age should be greater than 0")
	}

	if len(a.Name) <= 8 {
		return fmt.Errorf("name length at least should be 8")
	}

	return nil
}
