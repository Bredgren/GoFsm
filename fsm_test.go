package fsm

import (
	"bufio"
	"os"
	"testing"
)

type state struct {
	name   string
	enters int
	exits  int
}

func (s *state) OnEnter() {
	s.enters++
}

func (s *state) OnExit() {
	s.exits++
}

func (s state) Name() string {
	return s.name
}

func TestFsm(t *testing.T) {
	s1 := &state{name: "s1"}
	s2 := &state{name: "s2"}
	var fsm *Fsm = NewFsm(s1)

	if s1.enters != 1 || s1.exits != 0 {
		t.Error("OnEnter not called on init:", s1.enters, s1.exits)
	}

	fsm.GotoState(s2)

	if s1.enters != 1 || s1.exits != 1 {
		t.Error("OnExit not called on switch:", s1.enters, s1.exits)
	}

	if s2.enters != 1 || s2.exits != 0 {
		t.Error("OnEnter not called on switch:", s2.enters, s2.exits)
	}

	if fsm.PreviousState.Name() != s1.Name() {
		t.Error("previousState incorrect:", fsm.PreviousState)
	}
}

// Testing with the example from http://en.wikipedia.org/wiki/Automata-based_programming
var char string

type beforeState struct {
	res *string
}

func (s *beforeState) Name() string {
	return "before"
}
func (s *beforeState) OnEnter() {}
func (s *beforeState) OnExit() {
	*s.res += char
}

type insideState struct {
	res *string
}

func (s *insideState) Name() string {
	return "inside"
}
func (s *insideState) OnEnter() {}
func (s *insideState) OnExit()  {}

type afterState struct {
	res *string
}

func (s *afterState) Name() string {
	return "after"
}
func (s *afterState) OnEnter() {}
func (s *afterState) OnExit() {
	*s.res += "\n"
}

func TestFsmExample(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanBytes)

	var res string
	before := &beforeState{&res}
	inside := &insideState{&res}
	after := &afterState{&res}
	var fsm *Fsm = NewFsm(before)
	for s.Scan() {
		char = s.Text()
		if char == "\n" {
			fsm.GotoState(before)
		} else {
			switch fsm.CurrentState.Name() {
			case "before":
				if char != " " {
					fsm.GotoState(inside)
				}
			case "inside":
				if char == " " {
					fsm.GotoState(after)
				} else {
					res += char
				}
			case "after":
			}
		}
	}

	exp := "abc\ndef\ng\n"
	if res != exp {
		t.Fatalf("Got: %s, Exp: %s", res, exp)
	}
}
