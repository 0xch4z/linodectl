package resource

import "fmt"

type NoPropertyError struct {
	Name string
}

func (e NoPropertyError) Error() string {
	return fmt.Sprintf("property %q does not exist", e.Name)
}

type NotUpdateableError struct {
	Name string
}

func (e NotUpdateableError) Error() string {
	return fmt.Sprintf("property %q is not updateable", e.Name)
}

type InvalidValueError struct {
	Name  string
	Value interface{}
}

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("invalid %s value: %v", e.Name, e.Value)
}
