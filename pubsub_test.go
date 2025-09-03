package util

import (
	"testing"
)

func TestNewPubSub(t *testing.T) {
	pubSub := NewPubSub[string]()

	one := pubSub.Register(1)
	two := pubSub.Register(1)
	three := pubSub.Register(1)

	pubSub.Broadcast("one")
	if <-one != "one" {
		t.Fail()
	}
	if <-two != "one" {
		t.Fail()
	}
	if <-three != "one" {
		t.Fail()
	}

	pubSub.Unregister(two)
	pubSub.Broadcast("one")
	if <-one != "one" {
		t.Fail()
	}
	if <-three != "one" {
		t.Fail()
	}

	pubSub.Unregister(one)
	pubSub.Broadcast("one")
	if <-three != "one" {
		t.Fail()
	}

	pubSub.Unregister(three)

	if pubSub.IsActive() {
		t.Fail()
	}
}
