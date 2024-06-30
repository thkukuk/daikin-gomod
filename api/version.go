package daikin

import (
        "strings"
)

// Version is the human-readable firmware version of the Daikin unit.
type Version struct {
	value string
	param string
}

func (v *Version) String() string {
	return v.value
}

func (v *Version) setUrlValues() string {
	return v.param + "=" + v.String()
}

func (v *Version) decode(param string, s string) error {
	// Replace "_" with "."
        str := strings.Replace(s, "_", ".", -1)
	*v = Version{value: str, param: param}
	return nil
}

