package daikin

import (
        "fmt"
        "net/url"
        "strconv"
)

// FanDir is the louvre swing setting of the Daikin unit.
type FanDir int

// Supported louve settings. Not all models will support all values.
const (
	FanDirStopped    FanDir = 0
	FanDirVertical   FanDir = 1
	FanDirHorizontal FanDir = 2
	FanDirBoth       FanDir = 3
)

var fanDirMap = map[FanDir]string{
	FanDirStopped:    "Stopped",
	FanDirVertical:   "Vertical",
	FanDirHorizontal: "Horizontal",
	FanDirBoth:       "Both",
}

func (f *FanDir) setUrlValues(v url.Values) {
	v.Set("f_dir", strconv.Itoa(int(*f)))
}

func (f *FanDir) decode(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid f_dir value: %s (err=%v)", s, err)
	}
	fd := FanDir(v)
	if _, ok := fanDirMap[fd]; !ok {
		return fmt.Errorf("unknown f_dir value: %s", s)
	}
	*f = fd
	return nil
}

func (f *FanDir) String() string {
	v, ok := fanDirMap[*f]
	if !ok {
		return fmt.Sprintf("Unknown FanDir [%d]", int(*f))
	}
	return v
}

func (f *FanDir) Float64() float64 {
	return float64(*f)
}

