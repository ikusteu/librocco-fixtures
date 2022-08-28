package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// #region StringSet
type StringSet struct {
	m     map[string]bool
	oKeys []string
}

func NewStringSet() *StringSet {
	m := make(map[string]bool)
	k := []string{}
	set := &StringSet{m, k}
	return set
}

func (s *StringSet) Add(input string) {
	if !s.Exists(input) {
		s.m[input] = true
		s.oKeys = append(s.oKeys, input)
	}
}

func (s *StringSet) Remove(input string) {
	if !s.Exists(input) {
		s.m[input] = false
		s.oKeys = append(s.oKeys, input)
	}
}

func (s *StringSet) Exists(input string) bool {
	_, ok := s.m[input]
	return ok
}

func (s *StringSet) AddSlice(slice []string) {
	for _, el := range slice {
		s.Add(el)
	}
}

func (s *StringSet) ToSlice() []string {
	return s.oKeys
}

func (s *StringSet) Print() {
	for _, item := range s.oKeys {
		fmt.Println(item)
	}
}

func (s *StringSet) GetRandom() string {
	i := RandInt(0, len(s.oKeys)-1)
	return s.oKeys[i]
}

// #region StringSet

func RandInt(start int, end int) int {
	rand.Seed(time.Now().UnixNano())

	rangeDiff := end - start

	rf := rand.Float32() * float32(rangeDiff+1)
	ri := int(rf) + start

	return ri
}
