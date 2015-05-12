package flags

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var Int = map[string]int{}
var String = map[string]string{}
var Duration = map[string]time.Duration{}

func Parse(s string) error {

	ints := map[string]*int{}
	strs := map[string]*string{}
	durns := map[string]*time.Duration{}

	lines := strings.Split(s, "\n")
	for ln, line := range lines {

		name, rest := field(line)
		if len(name) == 0 {
			continue
		}

		// comment
		if name[0] == '#' {
			continue
		}

		// error default value
		if len(rest) == 0 {
			return fmt.Errorf("line %d expecting default value", ln)
		}

		value, therest := field(rest)

		// error usage
		if len(therest) == 0 {
			return fmt.Errorf("line %d expecting usage text", ln)
		}

		// flag int
		intv, err := strconv.Atoi(value)
		if err == nil {
			ints[name] = flag.Int(name, intv, therest)
			continue
		}

		// flag duration
		durnv, err := time.ParseDuration(value)
		if err == nil {
			durns[name] = flag.Duration(name, durnv, therest)
			continue
		}

		// flag string
		strs[name] = flag.String(name, value, therest)
	}

	flag.Parse()

	// map int
	for name, value := range ints {
		Int[name] = *value
	}

	// map duration
	for name, value := range durns {
		Duration[name] = *value
	}

	// map string
	for name, value := range strs {
		String[name] = *value
	}

	return nil
}

func field(s string) (next, rest string) {

	var i int
	var c rune

	for i, c = range s {
		if !unicode.IsSpace(c) {
			break
		}
	}

	if i == len(s) {
		return
	}

	start := i
	for i, c = range s[start:] {
		if unicode.IsSpace(c) {
			break
		}
	}

	end := start + i
	next = s[start:end]
	if end == len(s) {
		return
	}

	for i, c = range s[end:] {
		if !unicode.IsSpace(c) {
			break
		}
	}
	rest = s[end+i:]
	return
}
