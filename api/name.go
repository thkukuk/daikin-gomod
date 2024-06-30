package daikin

import (
        "net/url"
)

// Name is the human-readable name of the Daikin unit.
type Name struct {
	value string
	param string
}

func (n *Name) String() string {
	return n.value
}

func (n *Name) setUrlValues() string {
	return n.param + "=" + url.PathEscape(n.String())
}

func (n *Name) decode(param string, s string) error {
	v, err := url.PathUnescape(s)
	if err != nil {
		return err
	}
	*n = Name{value: v, param: param}
	return nil
}

