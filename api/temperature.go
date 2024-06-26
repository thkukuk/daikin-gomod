package daikin

import (
        "fmt"
        "net/url"
        "strconv"
)

// Temperature in Celcius.
type Temperature struct {
	value float64
	param string
}

func (t *Temperature) setUrlValues(v url.Values) {
	v.Set(t.param, t.String())
}

func (t *Temperature) decode(param string, v string) error {

	if v == "--" {
	   v = "-1"
	} 

	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("Temperature: error parsing %s=%s: %v", param, v, err)
	}
	*t = Temperature{value: val, param: param}
	return nil
}

func (t *Temperature) String() string {
        if t.value == -1 {
                return "N/A"
        }
	return strconv.FormatFloat(t.value, 'f', 1, 64)
}

func (t *Temperature) Float64() float64 {
	return t.value
}
