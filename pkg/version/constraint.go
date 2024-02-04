package version

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

type operator uint8

const (
	equal operator = iota
	notEqual
	greaterThan
	lessThan
	greaterThanEqual
	lessThanEqual
	tilde
	caret
	glob
	anything
)

const separator = ","

var errInvalid = errors.New("invalid version constraint")

func parseOperator(v string) (operator, string, bool) {
	var i int
	for i < len(v) && (v[i] == ' ') {
		i++
	}

	if len(v) < i+1 {
		return anything, "", true // TODO: Decide if this should return an error.
	}

	switch v[i] {
	case '=':
		return equal, v[i+1:], true
	case '!':
		if len(v[i:]) > 1 && v[i+1] == '=' {
			return notEqual, v[i+2:], true
		}
		return 0, v, false
	case '>':
		if len(v[i:]) > 1 && v[i+1] == '=' {
			return greaterThanEqual, v[i+2:], true
		}
		return greaterThan, v[i+1:], true
	case '<':
		if len(v[i:]) > 1 && v[i+1] == '=' {
			return lessThanEqual, v[i+2:], true
		}
		return lessThan, v[i+1:], true
	case '~':
		if len(v[i:]) > 1 && (v[i+1] == '>' || v[i+1] == '=') {
			i++
		}
		return tilde, v[i+1:], true
	case '^':
		return caret, v[i+1:], true
	case '*':
		return anything, v[i+1:], true
	case 'v', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if e := strings.IndexByte(v, '*'); e != -1 {
			return glob, v[i:e], true
		}
		return equal, v[i:], true
	}

	return 0, v, false
}

// Constraint represents a version constraint.
type Constraint []constraint

// NewConstraint returns a new version constraint.
func NewConstraint(v string) (Constraint, error) {
	switch v {
	case "", "*", "latest":
		return nil, nil
	}

	res := make([]constraint, 0, strings.Count(v, separator)+1)

	for v != "" {
		var (
			op operator
			c  string
			ok bool
		)
		c, v, _ = strings.Cut(v, separator)
		if c == "" {
			continue // TODO: Decide if this should return an error.
		}
		op, c, ok = parseOperator(c)
		if !ok {
			return nil, fmt.Errorf("%w: %s", errInvalid, c)
		}

		if op != anything {
			c = clean(c)
			var valid bool
			if op == glob {
				valid = semver.IsValid(strings.TrimSuffix(c, "."))
			} else {
				valid = semver.IsValid(c)
			}
			if !valid {
				return nil, fmt.Errorf("%w: %s", errInvalid, c)
			}
			res = append(res, constraint{op: op, v: c})
		} else if c != "" {
			return nil, fmt.Errorf("%w: *%s", errInvalid, c)
		}
	}

	return res, nil
}

// Check returns true if the version satisfies the constraint.
func (c Constraint) Check(v string) bool {
	v = clean(v)
	for _, c := range c {
		if !c.check(v) {
			return false
		}
	}
	return true
}

type constraint struct {
	op operator
	v  string
}

func (c constraint) check(v string) bool {
	switch c.op {
	case equal:
		return semver.Compare(c.v, v) == 0
	case notEqual:
		return semver.Compare(c.v, v) != 0
	case greaterThan:
		return semver.Compare(v, c.v) > 0
	case lessThan:
		return semver.Compare(v, c.v) < 0
	case greaterThanEqual:
		return semver.Compare(v, c.v) >= 0
	case lessThanEqual:
		return semver.Compare(v, c.v) <= 0
	case tilde:
		if semver.Compare(v, c.v) < 0 {
			return false
		}
		if strings.IndexByte(c.v, '.') == -1 {
			return strings.HasPrefix(v, c.v)
		}
		return semver.MajorMinor(v) == semver.MajorMinor(c.v)
	case caret:
		if semver.Compare(v, c.v) < 0 {
			return false
		}
		if strings.HasPrefix(c.v, "v0.") {
			return semver.MajorMinor(v) == semver.MajorMinor(c.v)
		}
		return semver.Major(v) == semver.Major(c.v)
	case glob:
		return strings.HasPrefix(v, c.v)
	default:
		return false // Does not happen.
	}
}

// Check returns true if the version satisfies the constraint.
func Check(constraint, version string) bool {
	c, err := NewConstraint(constraint)
	if err != nil {
		return false
	}
	return c.Check(version)
}

func clean(v string) string {
	i := 0
	for i < len(v) && (v[i] == ' ') {
		i++
	}
	n := len(v)
	for n > i && (v[n-1] == ' ') {
		n--
	}
	if v[i] != 'v' {
		return "v" + v[i:n]
	}
	return v[i:n]
}
