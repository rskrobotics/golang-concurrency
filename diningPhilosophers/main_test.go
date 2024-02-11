package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		order.order = []string{}
		dine()
		if len(order.order) != len(philosophers) {
			t.Errorf("Incorrect number of finished slice. Expected %d, but got %d\n", len(philosophers), len(order.order))
		}
	}
}

func Test_dineWithVaryingDelays(t *testing.T) {
	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", time.Second * 0},
		{"quater second delay", time.Millisecond * 250},
		{"half second delay", time.Millisecond * 500},
	}

	for _, e := range theTests {
		order.order = []string{}
		eatTime = e.delay
		thinkTime = e.delay
		dine()
		if len(order.order) != len(philosophers) {
			t.Errorf("%s, Incorrect number of finished slice. Expected %d, but got %d\n", e.name, len(philosophers), len(order.order))
		}
	}
}
