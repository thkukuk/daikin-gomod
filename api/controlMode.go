package daikin

import (
        "fmt"
        "net/url"
        "strconv"
)

// Mode is the operating mode of the Daikin unit.
type Mode int

// The valid modes supported by the Daikin Wifi module (not all units
// may support all values).
const (
	ModeDehumidify Mode = 2
	ModeCool       Mode = 3
	ModeHeat       Mode = 4
	ModeFan        Mode = 6
	ModeAuto       Mode = 0
	ModeAuto1      Mode = 1
	ModeAuto7      Mode = 7
)

var modeMap = map[Mode]string{
	ModeDehumidify: "Dehumidify",
	ModeCool:       "Cool",
	ModeHeat:       "Heat",
	ModeFan:        "Fan",
	ModeAuto:       "Auto",
	ModeAuto1:      "Auto",
	ModeAuto7:      "Auto",
}

func (m *Mode) setUrlValues(v url.Values) {
	v.Set("mode", strconv.Itoa(int(*m)))
}

func (m *Mode) decode(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid mode value: %s (err=%v)", s, err)
	}
	i := Mode(v)
	if _, ok := modeMap[i]; !ok {
		return fmt.Errorf("unknown mode value: %s", s)
	}
	*m = i

	return nil
}

func (m *Mode) String() string {
	if v, ok := modeMap[*m]; ok {
		return v
	}
	return fmt.Sprintf("Unknown Mode [%d]", *m)
}

func (m *Mode) Float64() float64 {
	return float64(*m)
}

