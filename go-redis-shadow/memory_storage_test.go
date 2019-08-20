package main

import (
	"testing"
)

func Test1(t *testing.T) {
	s := NewMemoryStorage()
	s.SetString("name", "Leo")
	v, ok := s.GetString("name")
	if !ok || v != "Leo" {
		t.Errorf("GetString = (%v, %v); want (Leo, true)", v, ok)
	}
	s.DeleteKey("name")
	_, ok = s.GetString("name")
	if ok {
		t.Errorf("DeleteKey = %v; want false", ok)
	}
}
