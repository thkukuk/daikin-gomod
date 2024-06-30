package daikin

import (
)

// String is a generic class for string values
type String struct {
	value string
	param string
}

func (s *String) String() string {
	return s.value
}

func (s *String) setUrlValues() string {
	return s.param + "=" + s.String()
}

func (s *String) decode(param string, str string) error {
	*s = String{value: str, param: param}
	return nil
}

