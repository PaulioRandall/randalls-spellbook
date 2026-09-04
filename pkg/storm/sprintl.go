package storm

import (
// "strings"
)

// Package sprintl

type Liner struct {
	// TODO
}

func NewLiner() *Liner {
	return &Liner{}
}

func (sl *Liner) String() string {
	// TODO
	return ""
}

func (sl *Liner) F(
	line int,
	args ...string,
) *Liner {
	// TODO
	return nil
}

func (sl *Liner) R(
	line int,
	values ...string,
) *Liner {
	// TODO
	return nil
}

func (sl *Liner) RF(
	line int,
	values ...[]string,
) *Liner {
	// TODO
	return nil
}

func (sl *Liner) J(
	line int,
	delim string,
	values ...string,
) *Liner {
	// TODO
	return nil
}

func (sl *Liner) JF(
	line int,
	delim string,
	values ...[]string,
) *Liner {
	// TODO
	return nil
}

func (sl *Liner) G(
	line int,
	generator func(index int, line string) string,
) *Liner {
	// TODO
	return nil
}

func F(
	s string,
	args ...string,
) string {
	// TODO
	return ""
}

func R(
	s string,
	values ...string,
) string {
	// TODO
	return ""
}

func RF(
	s string,
	values ...[]string,
) string {
	// TODO
	return ""
}

func J(
	s string,
	delim string,
	values ...string,
) string {
	// TODO
	return ""
}

func JF(
	s string,
	delim string,
	values ...[]string,
) string {
	// TODO
	return ""
}

func G(
	s string,
	generator func(index int, line string) string,
) string {
	// TODO
	return ""
}
