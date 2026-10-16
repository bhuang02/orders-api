package main

import "testing"

func TestBasicMath(t *testing.T) {
    if 2+2 != 4 { t.Fatal("unexpected result") }
}
