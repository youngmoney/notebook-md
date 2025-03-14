package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type Range struct {
	Lower int
	Upper int
}

func RangeLowerUpper(l, u int) Range {
	return Range{Lower: l, Upper: u}
}

func RangeOpenUpper(u int) Range {
	return Range{Lower: 0, Upper: u}
}

func RangeLowerOpen(l int) Range {
	return Range{Lower: l, Upper: math.MaxInt}
}

func RangeOpenOpen() Range {
	return Range{Lower: 0, Upper: math.MaxInt}
}

func (r Range) Contains(i int) bool {
	return i >= r.Lower && i <= r.Upper
}

func (r Range) Overlaps(o Range) bool {
	if o.Upper < r.Lower || o.Lower > r.Upper {
		return false
	}
	return true
}

type MultiRange struct {
	Ranges []Range
}

func (m *MultiRange) Add(r Range) {
	m.Ranges = append(m.Ranges, r)
}

func (m MultiRange) Contains(i int) bool {
	for _, r := range m.Ranges {
		if r.Contains(i) {
			return true
		}
	}
	return false
}

func (m MultiRange) Overlaps(o Range) bool {
	for _, r := range m.Ranges {
		if r.Overlaps(o) {
			return true
		}
	}
	return false
}

func asNumbers(s []string) ([]int, error) {
	var i []int
	for _, p := range s {
		if p == "" {
			i = append(i, -1)
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return []int{}, errors.New("invalid line number: " + p)
		}
		i = append(i, n)
	}
	return i, nil

}

func MultiRangeFromString(s string) (MultiRange, error) {
	var m MultiRange
	if s == "" {
		m.Add(RangeOpenOpen())
		return m, nil
	}
	for _, l := range strings.Split(s, ",") {
		p := strings.Split(l, "-")
		n, err := asNumbers(p)
		if err != nil {
			return m, err
		}
		switch len(p) {
		case 0:
			continue
		case 1:
			m.Add(RangeLowerUpper(n[0], n[0]))
		case 2:
			if n[0] == -1 && n[1] == -1 {
				m.Add(RangeOpenOpen())
			} else if n[0] == -1 {
				m.Add(RangeOpenUpper(n[1]))
			} else if n[1] == -1 {
				m.Add(RangeLowerOpen(n[0]))
			} else {
				m.Add(RangeLowerUpper(n[0], n[1]))
			}
		default:
			return m, errors.New("unknow range: " + l)
		}
	}
	return m, nil
}
