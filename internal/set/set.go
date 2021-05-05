package set

import (
	"fmt"
	"sort"
)

type (
	StringSet interface {
		Add(x ...string)
		Del(x ...string)
		In(x ...string) bool
		Len() int
		Unwrap() []string
		Equal(x StringSet) bool
		And(x StringSet) StringSet
		Clone() StringSet
	}

	stringSet struct {
		v map[string]bool
	}
)

func NewStringSet(seed ...string) StringSet {
	s := &stringSet{
		v: map[string]bool{},
	}
	s.Add(seed...)
	return s
}

func (s stringSet) And(x StringSet) StringSet {
	p := NewStringSet()
	for _, e := range x.Unwrap() {
		if s.In(e) {
			p.Add(e)
		}
	}
	return p
}

func (s stringSet) Equal(x StringSet) bool {
	if x == nil {
		return false
	}
	var (
		xs = s.Unwrap()
		ys = x.Unwrap()
	)
	if len(xs) != len(ys) {
		return false
	}
	sort.Strings(xs)
	sort.Strings(ys)
	for i, p := range xs {
		if p != ys[i] {
			return false
		}
	}
	return true
}
func (s *stringSet) Add(x ...string) {
	for _, p := range x {
		s.v[p] = true
	}
}
func (s *stringSet) Del(x ...string) {
	for _, p := range x {
		delete(s.v, p)
	}
}
func (s stringSet) In(x ...string) bool {
	for _, p := range x {
		if !s.v[p] {
			return false
		}
	}
	return true
}
func (s stringSet) Len() int { return len(s.v) }
func (s stringSet) Unwrap() []string {
	var (
		i int
		v = make([]string, len(s.v))
	)
	for k := range s.v {
		v[i] = k
		i++
	}
	return v
}
func (s stringSet) Clone() StringSet { return NewStringSet(s.Unwrap()...) }
func (s stringSet) String() string   { return fmt.Sprint(s.Unwrap()) }

type (
	RuneSet interface {
		Add(x ...rune)
		Del(x ...rune)
		In(x ...rune) bool
		Len() int
		Unwrap() []rune
		Equal(x RuneSet) bool
		Clone() RuneSet
	}

	runeSet struct {
		v map[rune]bool
	}
)

func NewRuneSet(seed ...rune) RuneSet {
	s := &runeSet{
		v: map[rune]bool{},
	}
	s.Add(seed...)
	return s
}

func (s *runeSet) Equal(x RuneSet) bool {
	if x == nil {
		return false
	}
	var (
		xs = s.Unwrap()
		ys = x.Unwrap()
	)
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	for i, p := range xs {
		if p != ys[i] {
			return false
		}
	}
	return true
}
func (s *runeSet) Add(x ...rune) {
	for _, p := range x {
		s.v[p] = true
	}
}
func (s *runeSet) Del(x ...rune) {
	for _, p := range x {
		delete(s.v, p)
	}
}
func (s runeSet) In(x ...rune) bool {
	for _, p := range x {
		if !s.v[p] {
			return false
		}
	}
	return true
}
func (s runeSet) Len() int { return len(s.v) }
func (s runeSet) Unwrap() []rune {
	var (
		i int
		v = make([]rune, len(s.v))
	)
	for k := range s.v {
		v[i] = k
		i++
	}
	return v
}
func (s runeSet) Clone() RuneSet { return NewRuneSet(s.Unwrap()...) }
func (s runeSet) String() string { return fmt.Sprint(s.Unwrap()) }
