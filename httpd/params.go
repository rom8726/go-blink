package httpd

import "strconv"

type Params map[string]string

func (p Params) Int(name string) int {
	v, ok := p[name]
	if !ok {
		return 0
	}

	i, _ := strconv.ParseInt(v, 10, 64)
	return int(i)
}

func (p Params) Int32(name string) int32 {
	v, ok := p[name]
	if !ok {
		return 0
	}

	i, _ := strconv.ParseInt(v, 10, 32)
	return int32(i)
}

func (p Params) Int64(name string) int64 {
	v, ok := p[name]
	if !ok {
		return 0
	}

	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}
