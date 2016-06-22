package models

import "testing"

func TestEncrytPass(t *testing.T) {
	en := EncrytPass("test")
	t.Log(en)
}
